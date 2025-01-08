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
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/config"
)

func NewDatabaseConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Running migrations
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Unable to create driver", "error", err)
		return nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		slog.Error("Unable to get working directory", "error", err)
		return nil, err
	}

	migrationPath := filepath.Join(basePath, "database", "migrations")

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "postgres", driver)
	if err != nil {
		slog.Error("Unable to create migration instance", "error", err)
		return nil, err
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			slog.Info("No migration to run")
		} else {
			slog.Error("Unable to run migration", "error", err)
			return nil, err
		}
	}

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
