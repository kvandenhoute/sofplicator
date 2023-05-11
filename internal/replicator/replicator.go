package replicator

import (
	"errors"

	"github.com/kvandenhoute/sofplicator/internal/config"
	"github.com/kvandenhoute/sofplicator/internal/util"
	batchv1 "k8s.io/api/batch/v1"
)

func StartReplication(replication Replication) (string, error) {
	uuid := util.GenerateUuid()

	targetUsername, targetPassword, err := getCredentialsForRegistry(replication.Target)
	if err != nil {
		return "", err
	}
	sourceUsername, sourcePassword, err := getCredentialsForRegistry(replication.Source)
	if err != nil {
		return "", err
	}
	secretName, err := CreateSecret(replication.Identifier, uuid, targetUsername, targetPassword, sourceUsername, sourcePassword, getCurrentNamespace())
	if err != nil {
		return "", err
	}
	configMapName, err := CreateConfigmap(replication.Identifier, uuid, replication.Images, replication.Charts, getCurrentNamespace())
	if err != nil {
		return "", err
	}
	jobImage := config.Get().JobImage.Registry + "/" + config.Get().JobImage.Repository + ":" + config.Get().JobImage.Tag
	_, err = CreateJob(replication.Identifier, uuid, jobImage, getCurrentNamespace(), config.Get().DockerCredentialsSecret, config.Get().MountPath, configMapName, secretName, replication.Source.Url, replication.Target.Url)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func getCredentialsForRegistry(registry Registry) (string, string, error) {
	var username string
	var password string
	if registry.UseCredentialsFromAzureVault {
		username = util.GetAzSecret(config.Get().AcrInfo.UsernameKey, registry.VaultURI)
		password = util.GetAzSecret(config.Get().AcrInfo.PasswordKey, registry.VaultURI)
		if (len(password) == 0 && len(username) != 0) || (len(username) == 0 && len(password) != 0) {
			return "", "", errors.New("please make sure both username and password are available in the Azure Vault. Or set useCredentialsFromAzureVault to false")
		}
	} else if registry.UseExistingSecret {

	} else {
		username = registry.Username
		password = registry.Password
		if (len(password) == 0 && len(username) != 0) || (len(username) == 0 && len(password) != 0) {
			return "", "", errors.New("please provide both registry.password and registry.username, or neither")
		}
	}
	return username, password, nil
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
