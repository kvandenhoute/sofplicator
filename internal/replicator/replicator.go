package replicator

import (
	"fmt"
	"strings"

	"github.com/kvandenhoute/sofplicator/internal/config"
	"github.com/kvandenhoute/sofplicator/internal/util"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
)

func StartReplication(replication Replication) (string, error) {
	splittedTarget := strings.Split(replication.TargetRegistry, ".")[0]
	name := fmt.Sprintf("%s-%s-repl", splittedTarget, replication.ReplicationType)
	username := util.GetAzSecret(config.Get().AcrInfo.UsernameKey, replication.VaultURI)
	password := util.GetAzSecret(config.Get().AcrInfo.PasswordKey, replication.VaultURI)
	uuid := util.GenerateUuid()
	secretName, err := CreateSecret(name, uuid, username, password, getCurrentNamespace())
	if err != nil {
		return "", err
	}
	configMapName, err := CreateConfigmap(name, uuid, replication.Images, replication.Charts, getCurrentNamespace())
	if err != nil {
		return "", err
	}
	_, err = CreateJob(name, uuid, config.Get().JobImage.Registry+"/"+config.Get().JobImage.Repository+"/"+config.Get().JobImage.Tag, getCurrentNamespace(), config.Get().DockerCredentialsSecret, config.Get().MountPath, configMapName, secretName, replication.TargetRegistry)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func StartGlobalReplication(replication Replication) (map[string]*string, error) {
	var jobIds map[string]*string = make(map[string]*string)
	registriesWithVault := util.GetAllACRsWithLabel(util.ListSubscriptions(), config.Get().AcrInfo.TargetLabelKey, config.Get().AcrInfo.TargetLabelValue)

	for _, registryWithVault := range registriesWithVault {
		log.Info("Start replication to  %+v", replication.TargetRegistry)
		replication.VaultURI = "https://" + *registryWithVault.Vault.Name + ".vault.azure.net/"
		replication.TargetRegistry = *registryWithVault.Registry.LoginServer
		jobId, err := StartReplication(replication)
		if err != nil {
			return nil, err
		}
		jobIds[replication.TargetRegistry] = &jobId
	}

	return jobIds, nil
}

func CleanReplication(uuid string) error {
	err := CleanupResources(uuid, getCurrentNamespace())
	if err != nil {
		return err
	}
	return nil
}

func GetJobStatus(uuid string) (batchv1.JobStatus, error) {
	job, err := GetJobOnLabel(uuid, getCurrentNamespace())
	if err != nil {
		return batchv1.JobStatus{}, err
	}
	return job.Status, nil
}
