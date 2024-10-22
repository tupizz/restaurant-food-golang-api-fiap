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

type ClientRepository interface {
	Create(ctx context.Context, client entity.Client) (entity.Client, error)
	GetByCpf(ctx context.Context, cpf string) (entity.Client, error)
}
type ProductFilter struct {
	Category string
	Page     int
	PageSize int
}

type OrderFilter struct {
	Page     int
	PageSize int
}

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetById(ctx context.Context, id int) (entity.Product, error)
	GetAll(ctx context.Context, filter *ProductFilter) ([]entity.Product, int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	GetByIds(ctx context.Context, ids []int) ([]entity.Product, int, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
	GetByID(ctx context.Context, id int) (entity.Order, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, filter *OrderFilter) ([]fiapRestaurantDb.GetAllOrdersRow, error)
}

type PaymentRepository interface {
	UpdateOrderPaymentStatus(ctx context.Context, orderId int, status entity.PaymentStatus) error
}

type PaymentTaxSettingsRepository interface {
	GetAll(ctx context.Context) ([]entity.PaymentTaxSettings, error)
}
