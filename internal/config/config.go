package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	config := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return config
}
