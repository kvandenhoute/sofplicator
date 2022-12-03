package replication

import "github.com/Netflix/go-env"

type ReplicationInfo struct {
	Source struct {
		Url      *string `env:"SOURCE_URL"`
		Username *string `env:"SOURCE_USERNAME"`
		Password *string `env:"SOURCE_PASSWORD"`
	}
	Target struct {
		Url      *string `env:"TARGET_URL" `
		Username *string `env:"TARGET_USERNAME"`
		Password *string `env:"TARGET_PASSWORD"`
	}

	Location struct {
		Images *string `env:"IMAGES_LOCATION,default=/input/images.json"`
		Charts *string `env:"CHARTS_LOCATION,default=/input/charts.json"`
	}

	ContinueOnError bool `env:"CONTINUE_ON_ERROR,default=false"`

	Extras env.EnvSet
}
