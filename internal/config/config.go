package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Redis struct {
	URL      string
	Password string
}

type Config struct {
	DatabaseURL string
	Redis       Redis
}

func LoadConfig() *Config {
	slog.Info("Loading config")
	viper.AutomaticEnv()

	slog.Info("DATABASE_URL", "value", viper.GetString("DATABASE_URL"))
	slog.Info("REDIS_URL", "value", viper.GetString("REDIS_URL"))
	slog.Info("REDIS_PASSWORD", "value", viper.GetString("REDIS_PASSWORD"))

	config := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		Redis: Redis{
			URL:      viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
	}

	if config.DatabaseURL == "" {
		slog.Error("DATABASE_URL is not set")
	}

	if config.Redis.URL == "" {
		slog.Error("REDIS_URL is not set")
	}

	return config
}
