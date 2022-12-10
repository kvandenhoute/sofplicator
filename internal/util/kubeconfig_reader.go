package util

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

var kubeconfigLock = &sync.Mutex{}
var kubeclientLock = &sync.Mutex{}

var (
	kubeconfig *rest.Config
	kubeclient *kubernetes.Clientset
)

func KubeConfig() *rest.Config {
	kubeconfigLock.Lock()
	defer kubeconfigLock.Unlock()

	if kubeconfig == nil {
		log.Info("Initial load of kubeconfig")
		kubeconfig = readKubeConfig()
	}
	return kubeconfig
}

func KubeClient() *kubernetes.Clientset {
	kubeclientLock.Lock()
	defer kubeclientLock.Unlock()

	if kubeclient == nil {
		log.Info("Initial load of kubeclient")
		kubeclient = getKubeClient()
	}
	return kubeclient
}

func readKubeConfig() *rest.Config {
	log.Debug("Reading in-cluster kubeconfig")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Debug("Failed reading in-cluster config, trying local kubeconfig")
		config = ctrl.GetConfigOrDie()
	}
	return config
}

func getKubeClient() *kubernetes.Clientset {
	config := KubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Failed to create K8s clientset")
	}

	return clientset
}
