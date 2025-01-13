package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
)

type PaymentRepository interface {
	UpdateOrderPaymentStatus(ctx context.Context, externalReference string, paymentMethod string, status domain.PaymentStatus) error
}
