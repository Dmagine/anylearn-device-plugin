package deviceplugin

import (
	"github.com/NVIDIA/go-gpuallocator/gpuallocator"
	kubeletClient "github.com/dmagine/anylearn-device-plugin/pkg/kubelet"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	log "github.com/sirupsen/logrus"
)

type AnylearnDevicePluginController struct {
	clientset     *kubernetes.Clientset
	kubeletClient *kubeletClient.KubeletClient
	stopCh        chan interface{}

	cachedDevices []*Device

	guaranteeDevicePlugin  *AnylearnDevicePlugin
	besteffortDeviceplugin *AnylearnDevicePlugin
}

func NewAnylearnDevicePluginController() (*AnylearnDevicePluginController, error) {
	controller := &AnylearnDevicePluginController{}
	controller.guaranteeDevicePlugin = &AnylearnDevicePlugin{
		controller:     controller,
		resourceName:   utils.GuaranteeGPU,
		socket:         pluginapi.DevicePluginPath + "anylearn-guarantee-gpu.sock",
		allocatePolicy: gpuallocator.NewBestEffortPolicy(),
	}
	controller.besteffortDeviceplugin = &AnylearnDevicePlugin{
		controller:     controller,
		resourceName:   utils.BestEffortGPU,
		socket:         pluginapi.DevicePluginPath + "anylearn-besteffort-gpu.sock",
		allocatePolicy: gpuallocator.NewBestEffortPolicy(),
	}
	return controller, nil
}

func (controller *AnylearnDevicePluginController) Start() error {
	controller.stopCh = make(chan interface{})
	if err := controller.InitK8SClientsetInCluster(); err != nil {
		return err
	}
	if err := controller.InitKubeletClientInCluster(); err != nil {
		return err
	}
	controller.cachedDevices = controller.Devices()
	go CheckHealth(controller.stopCh, controller.cachedDevices, []chan<- *Device{
		controller.guaranteeDevicePlugin.healthCh,
		controller.besteffortDeviceplugin.healthCh,
	})
	if err := controller.guaranteeDevicePlugin.Start(); err != nil {
		return err
	}
	if err := controller.besteffortDeviceplugin.Start(); err != nil {
		return err
	}
	return nil
}

func (controller *AnylearnDevicePluginController) Stop() error {
	close(controller.stopCh)
	controller.stopCh = nil
	controller.clientset = nil
	controller.kubeletClient = nil
	controller.cachedDevices = nil
	if err := controller.guaranteeDevicePlugin.Stop(); err != nil {
		return err
	}
	if err := controller.besteffortDeviceplugin.Stop(); err != nil {
		return err
	}
	return nil
}

func (controller *AnylearnDevicePluginController) GetPlugins() []*AnylearnDevicePlugin {
	return []*AnylearnDevicePlugin{
		controller.guaranteeDevicePlugin,
		controller.besteffortDeviceplugin,
	}
}

func (controller *AnylearnDevicePluginController) InitK8SClientsetInCluster() error {
	var err error
	var config *rest.Config

	config, err = rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed due to %v", err)
	}

	controller.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed due to %v", err)
	}

	return nil
}

func (controller *AnylearnDevicePluginController) InitKubeletClientInCluster() error {
	var err error
	controller.kubeletClient, err = kubeletClient.NewKubeletClientInCluster()
	if err != nil {
		return err
	}
	return nil
}
