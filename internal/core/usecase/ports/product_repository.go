package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
)

type ProductFilter struct {
	Category string
	Page     int
	PageSize int
}

type ProductRepository interface {
	GetAll(ctx context.Context, filter *ProductFilter) ([]entities.Product, int, error)
	Create(ctx context.Context, product entities.Product) (entities.Product, error)
	GetById(ctx context.Context, id int) (entities.Product, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product entities.Product) (entities.Product, error)
	GetByIds(ctx context.Context, ids []int) ([]entities.Product, int, error)
}
