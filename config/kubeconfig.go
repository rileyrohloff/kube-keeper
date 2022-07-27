package config

import (
	"flag"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func K8sClient(env string) *kubernetes.Clientset {
	if env == "local" || env == "" {
		log.Info().Msg("loading in cluster config")
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Err(err).Msg("failed to create in-cluster config")
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Error().Err(err).Msg("failed to create k8s config")
		}
		return client
	} else if env == "remote" {
		log.Info().Msg("attempting to build remote client k8s config")
		client, err := buildRemoteClientConfig()
		if err != nil {
			log.Err(err).Msg("failed to build remote client config")
		}
		return client
	}
	return nil
}

func buildRemoteClientConfig() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	return client, nil
}
