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

	// Set connection pool parameters if necessary
	dbpool.Config().MaxConns = 50                     // Adjust the number of maximum connections
	dbpool.Config().MaxConnIdleTime = 5 * time.Minute // Set idle timeout for connections
	dbpool.Config().MaxConnLifetime = time.Hour       // Set maximum lifetime of a connection

	return dbpool, nil
}
