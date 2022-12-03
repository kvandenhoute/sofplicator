package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

func ReadKubeConfig() *rest.Config {
	log.Debug("Reading in-cluster kubeconfig")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Debug("Failed reading in-cluster config, trying local kubeconfig")
		config = ctrl.GetConfigOrDie()
	}
	return config
}

func GetKubeClient() *kubernetes.Clientset {
	config := ReadKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Failed to create K8s clientset")
	}

	return clientset
}
