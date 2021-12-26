package main

import (
	"os"
	"syscall"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/dmagine/anylearn-device-plugin/pkg/agent"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var sigs chan os.Signal
var anylearnAgent *agent.AnylearnAgent

func main() {
	c := cli.NewApp()
	c.Action = start
	c.Flags = []cli.Flag{}

	err := c.Run(os.Args)
	utils.FatalWhenError(err)
}

func start(c *cli.Context) (err error) {
	log.Info("Init K8S Components")
	err = InitK8SComponents()
	if err != nil {
		return
	}
	defer close(k8sStopCh)

	log.Info("Loading NVML")
	err = nvml.Init()
	if err != nil {
		return
	}
	defer func() { log.Info("Shutdown of NVML returned:", nvml.Shutdown()) }()

	log.Info("Init AnylearnAgent")
	anylearnAgent, err = agent.NewAnylearnAgent(clientset, kubeletClient, &podLister)
	if err != nil {
		return
	}

	log.Info("Starting OS watcher.")
	sigs = utils.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err != nil {
		return
	}
	if err = anylearnAgent.Start(); err != nil {
		return
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
	return
}
