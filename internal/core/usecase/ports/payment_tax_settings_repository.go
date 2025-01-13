package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
)

type PaymentTaxSettingsRepository interface {
	GetAll(ctx context.Context) ([]domain.PaymentTaxSettings, error)
}
