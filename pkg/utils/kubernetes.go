package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewK8SClientsetInCluster() (clientset *kubernetes.Clientset, err error) {
	var config *rest.Config
	config, err = rest.InClusterConfig()
	if err != nil {
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	return
}
