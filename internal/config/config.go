package config

import (
	"sync"

	"github.com/Netflix/go-env"
)

type JobImage struct {
	Registry   string `env:"JOB_REGISTRY,default:harbor.aks-we-devops-harbor.int.sofico.be"`
	Repository string `env:"JOB_REPOSITORY,default:dev/acr-skopeo-replicate-kvdh"`
	Tag        string `env:"JOB_TAG,default:1.0.0"`
}

type AcrInfo struct {
	UsernameKey      string `env:"ACR_USERNAME_KEY,default:acr-writer-username"`
	PasswordKey      string `env:"ACR_PASSWORD_KEY,default:acr-writer-password"`
	TargetLabelKey   string `env:"ACR_TARGET_LABEL_KEY,default:replicationTarget"`
	TargetLabelValue string `env:"ACR_TARGET_LABEL_VALUE,default:true"`
}

type ReplicationInfo struct {
	LogLevel                string `env:"LOG_LEVEL,default=trace"`
	DockerCredentialsSecret string `env:"DOCKER_CREDENTIALS_SECRET,default=docker-credentials"`
	MountPath               string `env:"DOCKER_CREDENTIALS_SECRET,default="/tmp/kvdh/sofplicator"`
	AcrInfo                 AcrInfo
	JobImage                JobImage
}

var lock = &sync.Mutex{}

var replicationInfo *ReplicationInfo

func Get() *ReplicationInfo {
	lock.Lock()
	defer lock.Unlock()

	if replicationInfo == nil {
		replicationInfo = &ReplicationInfo{}
		env.UnmarshalFromEnviron(&replicationInfo)
	}
	return replicationInfo
}
