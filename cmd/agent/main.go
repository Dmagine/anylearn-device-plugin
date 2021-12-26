/*
 * Copyright (c) 2019-2021, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"syscall"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	anyplugin "github.com/dmagine/anylearn-device-plugin/pkg/deviceplugin"
	anyprobe "github.com/dmagine/anylearn-device-plugin/pkg/sysprobe"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"github.com/fsnotify/fsnotify"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var devicePluginController *anyplugin.AnylearnDevicePluginController
var sysProbe *anyprobe.SysProbe
var watcher *fsnotify.Watcher
var sigs chan os.Signal

func init() {
	devicePluginController = &anyplugin.AnylearnDevicePluginController{}
}

func main() {
	c := cli.NewApp()
	c.Action = start

	c.Flags = []cli.Flag{}

	err := c.Run(os.Args)
	utils.FatalWhenError(err)
}

func start(c *cli.Context) error {
	var err error

	log.Info("Loading NVML")
	err = nvml.Init()
	if err != nil {
		return err
	}
	defer func() { log.Info("Shutdown of NVML returned:", nvml.Shutdown()) }()

	log.Info("Starting FS watcher.")
	watcher, err = utils.NewFSWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		return err
	}
	defer watcher.Close()

	log.Info("Starting OS watcher.")
	sigs = utils.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err != nil {
		return err
	}
restart:
	err = devicePluginController.Stop()
	if err != nil {
		return err
	}
	err = devicePluginController.Start()
	if err != nil {
		return err
	}

	log.Info("Init plugins.")
	devicePluginController, err = anyplugin.NewAnylearnDevicePluginController()
	if err != nil {
		return err
	}

events:
	// Start an infinite loop, waiting for several indicators to either log
	// some messages, trigger a restart of the plugins, or exit the program.
	for {
		select {
		// Detect a kubelet restart by watching for a newly created
		// 'pluginapi.KubeletSocket' file. When this occurs, restart this loop,
		// restarting all of the plugins in the process.
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				log.Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
				goto restart
			}

		// Watch for any other fs errors and log them.
		case err := <-watcher.Errors:
			log.Infof("inotify: %s", err)

		// Watch for any signals from the OS. On SIGHUP, restart this loop,
		// restarting all of the plugins in the process. On all other
		// signals, exit the loop and exit the program.
		case s := <-sigs:
			switch s {
			case syscall.SIGHUP:
				log.Info("Received SIGHUP, restarting.")
				goto restart
			default:
				log.Infof("Received signal \"%v\", shutting down.", s)
				devicePluginController.Stop()
				break events
			}
		}
	}
	return nil
}
