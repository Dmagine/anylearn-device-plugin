package main

import (
	"time"

	log "github.com/golang/glog"

	v1 "k8s.io/api/core/v1"
)

var (
	retries = 8
)

func (controller *AnylearnDeviceController) GetPodListWithRetry() (*v1.PodList, error) {
	podList, err := controller.GetPodList()
	for i := 0; i < retries && err != nil; i++ {
		podList, err = controller.GetPodList()
		log.Warningf("failed to get pending pod list, retry")
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		log.Warningf("not found from kubelet /pods api")
		return nil, err
	}
	return podList, nil
}

func (controller *AnylearnDeviceController) GetPodList() (*v1.PodList, error) {
	podList, err := controller.kubeletClient.GetNodeRunningPods()
	if err != nil {
		return nil, err
	}
	return podList, nil
}
