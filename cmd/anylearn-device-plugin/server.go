package main

import (
	"os"

	"github.com/dmagine/anylearn-device-plugin/pkg/kubelet/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/golang/glog"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Device struct {
	pluginapi.Device
	Paths []string
	Index string
}

type AnylearnDeviceController struct {
	clientset     *kubernetes.Clientset
	kubeletClient *client.KubeletClient

	guaranteeDevicePlugin  *AnylearnDevicePlugin
	besteffortDeviceplugin *AnylearnDevicePlugin
}

func (controller *AnylearnDeviceController) GetPlugins() []*AnylearnDevicePlugin {
	return []*AnylearnDevicePlugin{
		controller.guaranteeDevicePlugin,
		controller.besteffortDeviceplugin,
	}
}

func (controller *AnylearnDeviceController) InitK8S() error {
	kubeconfigFile := os.Getenv("KUBECONFIG")
	var err error
	var config *rest.Config

	if _, err = os.Stat(kubeconfigFile); err != nil {
		log.V(5).Infof("kubeconfig %s failed to find due to %v", kubeconfigFile, err)
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed due to %v", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		if err != nil {
			log.Fatalf("Failed due to %v", err)
		}
	}

	controller.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed due to %v", err)
	}

	// TODO: init Kubelet Client
	return nil
}

func NewAnylearnDeviceController() (*AnylearnDeviceController, error) {
	return nil, nil
}
