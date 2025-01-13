package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type PaymentTaxSettingsRepository interface {
	GetAll(ctx context.Context) ([]entities.PaymentTaxSettings, error)
}
