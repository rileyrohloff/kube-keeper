package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rileyrohloff/kube-keeper/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("failed to load .env file for config")
	}

	// default config values
	ENV := os.Getenv("RUN_ENV")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if ENV == "" {
		log.Info().Msg("failed to find RUN_ENV environment variable, defaulting to in-cluster config")
	}

	client := config.K8sClient(ENV)
	log.Info().Msg("successfully built k8s client")
	for {
		time.Sleep(time.Second * 5)
		pods, err := client.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("runtime error getting pods")
		}
		for _, value := range pods.Items {
			log.Info().Msg(fmt.Sprintf("pod: %v", value.Name))
		}
	}
}
