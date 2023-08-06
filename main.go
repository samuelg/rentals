package main

import (
	"fmt"
	"github.com/samuelg/rentals/config"
	log "github.com/samuelg/rentals/logging"
	"os"
)

func getEnv() string {
	if value, ok := os.LookupEnv("ENV"); ok {
		return value
	}
	// default to development
	return "development"
}

func main() {
	env := getEnv()
	config.Init(env)
	log.Init(config.GetConfig().LogLevel, config.GetConfig().AppVersion)

	log.Log.Info(fmt.Sprintf("Loaded config for %s environment", env))

	log.Log.Info("Rental API service")
}
