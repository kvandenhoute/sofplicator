package main

import (
	"github.com/kvandenhoute/sofplicator/api/api/v0/router"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func main() {
	var logLevel string

	log.Info("Starting up.. Reading arguments")
	flag.StringVar(&logLevel, "log-level", "trace", "Loglevel: trace, debug, info, warning, error, fatal, panic")
	flag.Parse()

	var level log.Level
	log.Info(logLevel)
	level.UnmarshalText([]byte(logLevel))
	log.SetLevel(level)

	router.ServeRouter()

}
