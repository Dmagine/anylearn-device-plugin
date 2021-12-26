package deviceplugin

import (
	"github.com/NVIDIA/go-gpuallocator/gpuallocator"
	kubeletClient "github.com/dmagine/anylearn-device-plugin/pkg/kubelet"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"k8s.io/client-go/kubernetes"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	log "github.com/sirupsen/logrus"
)

type AnylearnDevicePluginController struct {
	clientset     *kubernetes.Clientset
	kubeletClient *kubeletClient.KubeletClient
	databus       *utils.DataBus

	stopCh  chan interface{}
	watcher *fsnotify.Watcher

	GuaranteeDevicePlugin  *AnylearnDevicePlugin
	BesteffortDeviceplugin *AnylearnDevicePlugin

	cachedDevices []*Device
}

func NewAnylearnDevicePluginController(
	clientset *kubernetes.Clientset,
	kubeletClient *kubeletClient.KubeletClient) (*AnylearnDevicePluginController, error) {
	controller := &AnylearnDevicePluginController{
		clientset:     clientset,
		kubeletClient: kubeletClient,
	}
	controller.GuaranteeDevicePlugin = &AnylearnDevicePlugin{
		controller:     controller,
		resourceName:   utils.GuaranteeGPU,
		socket:         pluginapi.DevicePluginPath + utils.GuaranteeSocket,
		allocatePolicy: gpuallocator.NewBestEffortPolicy(),
	}
	controller.BesteffortDeviceplugin = &AnylearnDevicePlugin{
		controller:     controller,
		resourceName:   utils.BestEffortGPU,
		socket:         pluginapi.DevicePluginPath + utils.BestEffortSocket,
		allocatePolicy: gpuallocator.NewBestEffortPolicy(),
	}
	return controller, nil
}

func (controller *AnylearnDevicePluginController) Start() (err error) {
	controller.stopCh = make(chan interface{})
	controller.cachedDevices = controller.Devices()
	go controller.checkDeviceHealth()
	// TODO: Check Device Is Taken Or Released
	if err := controller.GuaranteeDevicePlugin.Start(); err != nil {
		return err
	}
	if err := controller.BesteffortDeviceplugin.Start(); err != nil {
		return err
	}
	log.Info("Starting FS watcher.")
	controller.watcher, err = utils.NewFSWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-controller.stopCh:
				return

			// Detect a kubelet restart by watching for a newly created
			// 'pluginapi.KubeletSocket' file. When this occurs, restart this loop,
			// restarting all of the plugins in the process.
			case event := <-controller.watcher.Events:
				if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
					log.Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
					controller.Restart()
				}

			// Watch for any other fs errors and log them.
			case err := <-controller.watcher.Errors:
				log.WithError(err).Error()
			}
		}
	}()
	return nil
}

func (controller *AnylearnDevicePluginController) Stop() error {
	var err error
	close(controller.stopCh)
	if err = controller.GuaranteeDevicePlugin.Stop(); err != nil {
		return err
	}
	if err = controller.BesteffortDeviceplugin.Stop(); err != nil {
		return err
	}
	if err = controller.watcher.Close(); err != nil {
		return err
	}
	controller.stopCh = nil
	controller.clientset = nil
	controller.kubeletClient = nil
	controller.cachedDevices = nil
	controller.watcher = nil
	return nil
}

func (controller *AnylearnDevicePluginController) Restart() error {
	var err error
	if err = controller.Stop(); err != nil {
		return err
	}
	if err = controller.Start(); err != nil {
		return err
	}
	return nil
}

func (controller *AnylearnDevicePluginController) GetPlugins() []*AnylearnDevicePlugin {
	return []*AnylearnDevicePlugin{
		controller.GuaranteeDevicePlugin,
		controller.BesteffortDeviceplugin,
	}
}
