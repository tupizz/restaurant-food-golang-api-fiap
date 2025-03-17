package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

const redisTimeout = 5 * time.Second

func NewRedisConnection(cfg *config.Config) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Testing connection.
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Unable to connect to Redis", "error", err)
		return nil, err
	}

	slog.Info("Connected to Redis successfully")

	return rdb, nil
}
