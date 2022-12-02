package replicator

import (
	"github.com/kvandenhoute/sofplicator/internal/util"
	log "github.com/sirupsen/logrus"
)

func StartReplication(replication Replication) string {
	username := util.GetAzSecret("acr-writer-username", replication.VaultURI)
	password := util.GetAzSecret("acr-writer-password", replication.VaultURI)
	secretName := CreateSecret(username, password, "dev-tooling")
	configMapName := CreateConfigmap(replication.Images, replication.Charts, "dev-tooling")
	jobId := CreateJob("harbor.aks-we-devops-harbor.int.sofico.be/dev/acr-skopeo-replicate-kvdh:1.0.0", "dev-tooling", "docker-credentials", "/tmp/kvdh/sofplicator", configMapName, secretName, replication.TargetRegistry)
	return jobId
}

func StartGlobalReplication(replication Replication) map[string]*string {
	var jobIds map[string]*string = make(map[string]*string)
	registriesWithVault := util.GetAllACRsWithLabel(util.ListSubscriptions(), "replicationTarget", "true")

	for _, registryWithVault := range registriesWithVault {
		log.Info("Start replication to  %+v", replication.TargetRegistry)
		replication.VaultURI = "https://" + *registryWithVault.Vault.Name + ".vault.azure.net/"
		replication.TargetRegistry = *registryWithVault.Registry.LoginServer
		jobId := StartReplication(replication)
		jobIds[replication.TargetRegistry] = &jobId
	}

	return jobIds
}
