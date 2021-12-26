package main

import (
	"time"

	"github.com/dmagine/anylearn-device-plugin/pkg/kubelet"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	informers "k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

var clientset *kubernetes.Clientset
var kubeletClient *kubelet.KubeletClient

var podListerSynced cache.InformerSynced
var podInformer coreinformers.PodInformer
var podLister corelisters.PodLister
var nodeStoreSynced cache.InformerSynced
var nodeInformer coreinformers.NodeInformer
var nodeLister corelisters.NodeLister
var deploymentSynced cache.InformerSynced
var deploymentInformer appsinformers.DeploymentInformer
var deploymentLister appslisters.DeploymentLister
var k8sStopCh chan struct{}

func InitK8SComponents() (err error) {
	k8sStopCh = make(chan struct{})
	clientset, err = utils.NewK8SClientsetInCluster()
	kubeletClient, err = kubelet.NewKubeletClientInCluster()

	informerFactory := informers.NewSharedInformerFactory(clientset, 60*time.Minute)
	podInformer = informerFactory.Core().V1().Pods()
	nodeInformer = informerFactory.Core().V1().Nodes()
	podListerSynced = podInformer.Informer().HasSynced
	nodeStoreSynced = nodeInformer.Informer().HasSynced
	/*
		deploymentInformer = informerFactory.Apps().V1().Deployments()
		deploymentSynced = deploymentInformer.Informer().HasSynced
	*/

	go informerFactory.Start(k8sStopCh)

	/*
		if !cache.WaitForCacheSync(k8sStopCh , nodeStoreSynced, podListerSynced, deploymentSynced) {
			log.Fatalln("未能与远端同步数据")
		}
	*/
	if !cache.WaitForCacheSync(k8sStopCh, nodeStoreSynced, podListerSynced) {
		log.Fatalln("未能与远端同步数据")
	}

	podLister = podInformer.Lister()
	nodeLister = nodeInformer.Lister()
	// deploymentLister = deploymentInformer.Lister()

	log.Infoln("K8S Informer已就绪")
	return
}
