package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	domainError "github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/error"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type healthCheckRepository struct {
	db *pgxpool.Pool
}

func NewHealthCheckRepository(db *pgxpool.Pool) ports.HealthCheckRepository {
	return &healthCheckRepository{db: db}
}

func (r *healthCheckRepository) Ping(ctx context.Context) error {
	if err := r.db.Ping(ctx); err != nil {
		return domainError.NewEntityNotProcessableError(err.Error(), "database")
	}

	return nil
}
