package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	slog.Info("Loading config")
	viper.AutomaticEnv()

	slog.Info("DATABASE_URL", "value", viper.GetString("DATABASE_URL"))

	config := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	if config.DatabaseURL == "" {
		slog.Error("DATABASE_URL is not set")
	}

	return config
}
