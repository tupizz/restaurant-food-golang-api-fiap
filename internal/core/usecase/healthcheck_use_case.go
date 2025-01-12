package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type HealthCheckUseCase interface {
	Ping(ctx context.Context) error
}

type healthCheckUseCase struct {
	healthCheckRepository ports.HealthCheckRepository
}

func NewHealthCheckUseCase(healthCheckRepository ports.HealthCheckRepository) HealthCheckUseCase {
	return &healthCheckUseCase{healthCheckRepository: healthCheckRepository}
}

func (h *healthCheckUseCase) Ping(ctx context.Context) error {
	return h.healthCheckRepository.Ping(ctx)
}
