package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	fmt.Println("Loading config")
	viper.AutomaticEnv()

	fmt.Println(viper.GetString("DATABASE_URL"))

	config := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return config
}
