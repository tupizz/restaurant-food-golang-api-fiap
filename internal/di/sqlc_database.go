package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

func NewSQLCDB(cfg *config.Config) (*fiapRestaurantDb.Queries, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		slog.Error("Unable to ping the database", "error", err)
		return nil, err
	}

	db := fiapRestaurantDb.New(dbpool)

	return db, nil
}
