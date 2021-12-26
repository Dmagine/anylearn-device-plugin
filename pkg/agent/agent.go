package agent

import (
	anyplugin "github.com/dmagine/anylearn-device-plugin/pkg/deviceplugin"
	kubeletClient "github.com/dmagine/anylearn-device-plugin/pkg/kubelet"
	anyprobe "github.com/dmagine/anylearn-device-plugin/pkg/sysprobe"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type AnylearnAgent struct {
	clientset     *kubernetes.Clientset
	kubeletClient *kubeletClient.KubeletClient
	databus       *utils.DataBus

	DevicePluginController *anyplugin.AnylearnDevicePluginController
	SysProbe               *anyprobe.SysProbe
}

func NewAnylearnAgent(
	clientset *kubernetes.Clientset,
	kubeletClient *kubeletClient.KubeletClient) (*AnylearnAgent, error) {
	return &AnylearnAgent{
		clientset:     clientset,
		kubeletClient: kubeletClient,
		databus:       utils.NewDataBus(),
	}, nil
}

func (agent *AnylearnAgent) Start() (err error) {
	log.Info("Init DevicePlugin")
	agent.DevicePluginController, err = anyplugin.NewAnylearnDevicePluginController(agent.clientset, agent.kubeletClient, agent.databus)
	if err != nil {
		return
	}
	err = agent.DevicePluginController.Start()
	if err != nil {
		return
	}
	return
}

func (agent *AnylearnAgent) Stop() (err error) {
	err = agent.DevicePluginController.Stop()
	if err != nil {
		return
	}
	agent.DevicePluginController = nil
	return
}

func (agent *AnylearnAgent) Restart() (err error) {
	err = agent.Stop()
	if err != nil {
		return
	}
	err = agent.Start()
	if err != nil {
		return
	}
	return
}
