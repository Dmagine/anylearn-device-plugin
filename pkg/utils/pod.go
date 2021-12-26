package utils

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func getGPUIDFromPodAnnotation(pod *v1.Pod) (uuid string) {
	if len(pod.ObjectMeta.Annotations) > 0 {
		if value, found := pod.ObjectMeta.Annotations[BestEffortGPUUUIDAnnotationKey]; found {
			return value
		}
	}
	return ""
}
func podIsNotRunning(pod v1.Pod) bool {
	status := pod.Status
	//deletionTimestamp
	if pod.DeletionTimestamp != nil {
		return true
	}

	// pod is scheduled but not initialized
	if status.Phase == v1.PodPending && podConditionTrueOnly(status.Conditions, v1.PodScheduled) {
		log.Infof("Pod %s only has PodScheduled, is not running", pod.Name)
		return true
	}

	return status.Phase == v1.PodFailed || status.Phase == v1.PodSucceeded || (pod.DeletionTimestamp != nil && notRunning(status.ContainerStatuses)) || (status.Phase == v1.PodPending && podConditionTrueOnly(status.Conditions, v1.PodScheduled))
}

// notRunning returns true if every status is terminated or waiting, or the status list
// is empty.
func notRunning(statuses []v1.ContainerStatus) bool {
	for _, status := range statuses {
		if status.State.Terminated == nil && status.State.Waiting == nil {
			return false
		}
	}
	return true
}

func podConditionTrue(conditions []v1.PodCondition, expect v1.PodConditionType) bool {
	for _, condition := range conditions {
		if condition.Type == expect && condition.Status == v1.ConditionTrue {
			return true
		}
	}

	return false
}

func podConditionTrueOnly(conditions []v1.PodCondition, expect v1.PodConditionType) bool {
	if len(conditions) != 1 {
		return false
	}

	for _, condition := range conditions {
		if condition.Type == expect && condition.Status == v1.ConditionTrue {
			return true
		}
	}

	return false
}
