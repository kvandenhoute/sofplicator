package main

import (
	"github.com/kvandenhoute/sofplicator/api/v0/router"
	"github.com/kvandenhoute/sofplicator/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	var level log.Level

	log.Info(config.Get().LogLevel)
	level.UnmarshalText([]byte(config.Get().LogLevel))
	log.SetLevel(level)

	router.ServeRouter()

}
