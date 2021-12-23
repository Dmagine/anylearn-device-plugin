package main

import (
	"fmt"
	"io/ioutil"
	"time"

	kubeletClient "github.com/dmagine/anylearn-device-plugin/pkg/kubelet/client"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
)

type AnylearnDeviceController struct {
	clientset     *kubernetes.Clientset
	kubeletClient *kubeletClient.KubeletClient

	cachedDevices []*Device

	guaranteeDevicePlugin  *AnylearnDevicePlugin
	besteffortDeviceplugin *AnylearnDevicePlugin
}

func (controller *AnylearnDeviceController) GetPlugins() []*AnylearnDevicePlugin {
	return []*AnylearnDevicePlugin{
		controller.guaranteeDevicePlugin,
		controller.besteffortDeviceplugin,
	}
}

func (controller *AnylearnDeviceController) InitK8SClientsetInCluster() error {
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

func (controller *AnylearnDeviceController) InitKubeletClientInCluster() error {
	var err error
	var token string
	var tokenByte []byte
	tokenByte, err = ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(fmt.Errorf("in cluster mode, find token failed, error: %v", err))
	}
	token = string(tokenByte)

	controller.kubeletClient, err = kubeletClient.NewKubeletClient(&kubeletClient.KubeletClientConfig{
		Address: "127.0.0.1",
		Port:    10250,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure:   true,
			ServerName: "kubelet",
		},
		BearerToken: token,
		HTTPTimeout: 20 * time.Second,
	})
	return nil
}

func NewAnylearnDeviceController() (*AnylearnDeviceController, error) {
	controller := &AnylearnDeviceController{}
	controller.InitK8SClientsetInCluster()
	controller.InitKubeletClientInCluster()
	controller.cachedDevices = controller.Devices()
	return nil, nil
}
