package ports

import (
	"context"

	sqlcDB "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
)

type OrderFilter struct {
	Page     int
	PageSize int
}

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
	GetByID(ctx context.Context, id int) (domain.Order, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, filter *OrderFilter) ([]sqlcDB.GetAllOrdersRow, error)
}
