package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type PaymentRepository interface {
	UpdateOrderPaymentStatus(ctx context.Context, externalReference string, paymentMethod string, status entities.PaymentStatus) error
}
