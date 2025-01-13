package domain

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

type OrderFilter struct {
	Page     int
	PageSize int
}

type PaymentRepository interface {
	UpdateOrderPaymentStatus(ctx context.Context, externalReference string, paymentMethod string, status entity.PaymentStatus) error
}
