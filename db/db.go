package db

import (
	"fmt"
	"github.com/samuelg/rentals/config"
	log "github.com/samuelg/rentals/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.GetConfig().DbHost,
		config.GetConfig().DbUser,
		config.GetConfig().DbPassword,
		config.GetConfig().DbName,
		config.GetConfig().DbPort,
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	log.Log.Info("Connected to database")
	DB = database
}
