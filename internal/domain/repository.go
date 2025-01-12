package domain

import (
	"context"

	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
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

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
	GetByID(ctx context.Context, id int) (entity.Order, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, filter *OrderFilter) ([]fiapRestaurantDb.GetAllOrdersRow, error)
}

type PaymentRepository interface {
	UpdateOrderPaymentStatus(ctx context.Context, externalReference string, paymentMethod string, status entity.PaymentStatus) error
}

type PaymentTaxSettingsRepository interface {
	GetAll(ctx context.Context) ([]entity.PaymentTaxSettings, error)
}
