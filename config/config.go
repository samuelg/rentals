package config

import (
	"github.com/spf13/viper"
	"log" // logrus won't be initialized yet
)

type Config struct {
	AppVersion      string `mapstructure:"app_version"`
	Host            string `mapstructure:"host"`
	Port            uint16 `mapstructure:"port"`
	LogLevel        string `mapstructure:"log_level"`
	DbHost          string `mapstructure:"db_host"`
	DbPort          uint16 `mapstructure:"db_port"`
	DbUser          string `mapstructure:"db_user"`
	DbPassword      string `mapstructure:"db_password"`
	DbName          string `mapstructure:"db_name"`
	DefaultApiLimit uint8  `mapstructure:"default_api_limit"`
}

var parsedConfig Config

func Init(env string) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(env)
	v.AddConfigPath("../config/")
	v.AddConfigPath("config/")
	v.AutomaticEnv()

	v.SetDefault("app_version", "1.0.0")
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 8080)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error parsing configuration file, %v", err)
	}
	if err := v.Unmarshal(&parsedConfig); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func GetConfig() *Config {
	return &parsedConfig
}
