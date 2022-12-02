package main

import (
	"github.com/kvandenhoute/sofplicator/internal/replicator"
	"github.com/kvandenhoute/sofplicator/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	var logLevel string = "INFO"

	var level log.Level
	level.UnmarshalText([]byte(logLevel))
	log.SetLevel(level)

	images := []replicator.Artifact{
		{
			Name:    "controller",
			Version: "v1.3.1",
			Repo:    "contrib",
		},
		{
			Name:    "kube-webhook-certgen",
			Version: "v1.3.0",
			Repo:    "contrib",
		},
		{
			Name:    "sealed-secrets-controller",
			Version: "v0.16.0",
			Repo:    "contrib",
		},
	}
	charts := []replicator.Artifact{
		{
			Name:    "ingress-nginx",
			Version: "4.2.5",
			Repo:    "contrib",
		},
		{
			Name:    "sealed-secrets",
			Version: "1.16.2",
			Repo:    "contrib",
		},
	}

	registries := util.GetAllACRsWithLabel(util.ListSubscriptions(), "test", "yes")

	for _, registry := range registries {
		configMapName := replicator.CreateConfigmap(images, charts, "default")
		replicator.CreateJob("harbor.aks-we-devops-harbor.int.sofico.be/dev/acr-skopeo-replicate-kvdh:1.0.0", "dev-tooling", "docker-credentials", "/etc/sofplicator", configMapName, "acr-credentials", *registry.Registry.LoginServer)
	}

}
