package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type HealthCheckPingUseCase interface {
	Run(ctx context.Context) error
}

type healthCheckPingUseCase struct {
	healthCheckRepository ports.HealthCheckRepository
}

func NewHealthCheckPingUseCase(healthCheckRepository ports.HealthCheckRepository) HealthCheckPingUseCase {
	return &healthCheckPingUseCase{healthCheckRepository: healthCheckRepository}
}

func (h *healthCheckPingUseCase) Run(ctx context.Context) error {
	return h.healthCheckRepository.Ping(ctx)
}
