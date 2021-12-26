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
	"github.com/dmagine/anylearn-device-plugin/pkg/agent"
	"github.com/dmagine/anylearn-device-plugin/pkg/kubelet"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var sigs chan os.Signal
var anylearnAgent *agent.AnylearnAgent

var clientset *kubernetes.Clientset
var kubeletClient *kubelet.KubeletClient

func init() {
	var err error
	clientset, err = utils.NewK8SClientsetInCluster()
	utils.FatalWhenError(err)
	kubeletClient, err = kubelet.NewKubeletClientInCluster()
	utils.FatalWhenError(err)
	anylearnAgent, err = agent.NewAnylearnAgent(clientset, kubeletClient)
	utils.FatalWhenError(err)
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

	log.Info("Starting OS watcher.")
	sigs = utils.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err != nil {
		return err
	}
	if err = anylearnAgent.Start(); err != nil {
		return err
	}
events:
	// Start an infinite loop, waiting for several indicators to either log
	// some messages, trigger a restart of the plugins, or exit the program.
	for {
		select {
		// Watch for any signals from the OS. On SIGHUP, restart this loop,
		// restarting all of the plugins in the process. On all other
		// signals, exit the loop and exit the program.
		case s := <-sigs:
			switch s {
			case syscall.SIGHUP:
				log.Info("Received SIGHUP, restarting.")
				anylearnAgent.Restart()
			default:
				log.Infof("Received signal \"%v\", shutting down.", s)
				anylearnAgent.Stop()
				break events
			}
		}
	}
	return nil
}
