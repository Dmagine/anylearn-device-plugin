package kubelet

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func (client *KubeletClient) GetNodeRunningPods() (*v1.PodList, error) {
	resp, err := client.client.Get(fmt.Sprintf("https://%v:%d/pods/", client.host, client.defaultPort))
	if err != nil {
		return nil, err
	}

	body, err := ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	podLists := &v1.PodList{}
	if err = json.Unmarshal(body, &podLists); err != nil {
		return nil, err
	}
	return podLists, err
}

func (client *KubeletClient) GetPodListWithRetry() (*v1.PodList, error) {
	podList, err := client.GetPodList()
	for i := 0; i < utils.RetryTimes && err != nil; i++ {
		podList, err = client.GetPodList()
		log.Warningf("failed to get pending pod list, retry")
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		log.Warningf("not found from kubelet /pods api")
		return nil, err
	}
	return podList, nil
}

func (client *KubeletClient) GetPodList() (*v1.PodList, error) {
	podList, err := client.GetNodeRunningPods()
	if err != nil {
		return nil, err
	}
	return podList, nil
}
