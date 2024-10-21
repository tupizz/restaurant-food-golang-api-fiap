package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

func NewDatabaseConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	// Optionally, ping the database to ensure connection is established
	err = dbpool.Ping(ctx)
	if err != nil {
		slog.Error("Unable to ping the database", "error", err)
		return nil, err
	}

	return dbpool, nil
}
