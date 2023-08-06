package main

import (
	"fmt"
	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
	"github.com/samuelg/rentals/server"
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

	db.Init()
	server.Init()
}
