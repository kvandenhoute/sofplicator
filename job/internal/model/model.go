package model

type Source struct {
	Url      string `env:"SOURCE_URL"`
	Username string `env:"SOURCE_USERNAME"`
	Password string `env:"SOURCE_PASSWORD"`
}

type Target struct {
	Url      string `env:"TARGET_URL" `
	Username string `env:"TARGET_USERNAME"`
	Password string `env:"TARGET_PASSWORD"`
}
type Location struct {
	Images string `env:"IMAGES_LOCATION,default=/app/input/images.json"`
	Charts string `env:"CHARTS_LOCATION,default=/app/input/charts.json"`
}
type ReplicationInfo struct {
	Source   Source
	Target   Target
	Location Location

	ContinueOnError bool `env:"CONTINUE_ON_ERROR,default=false"`
}

type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}
