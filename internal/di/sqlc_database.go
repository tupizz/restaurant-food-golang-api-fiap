package di

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

func NewSQLCDB(cfg *config.Config) (*fiapRestaurantDb.Queries, *pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Running migrations
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Unable to create driver", "error", err)
		return nil, nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		slog.Error("Unable to get working directory", "error", err)
		return nil, nil, err
	}

	migrationPath := filepath.Join(basePath, "database", "migrations")

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "postgres", driver)
	if err != nil {
		slog.Error("Unable to create migration instance", "error", err)
		return nil, nil, err
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			slog.Info("No migration to run")
		} else {
			slog.Error("Unable to run migration", "error", err)
			return nil, nil, err
		}
	}

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, nil, err
	}

	dbpool.Config().MaxConns = 50
	dbpool.Config().MaxConnIdleTime = 5 * time.Minute
	dbpool.Config().MaxConnLifetime = time.Hour

	dbQuerie := fiapRestaurantDb.New(dbpool)

	return dbQuerie, dbpool, nil
}
