package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

func NewSQLCDB(cfg *config.Config) (*fiapRestaurantDb.Queries, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	dbpool.Config().MaxConns = 50
	dbpool.Config().MaxConnIdleTime = 5 * time.Minute
	dbpool.Config().MaxConnLifetime = time.Hour

	db := fiapRestaurantDb.New(dbpool)

	return db, nil
}
