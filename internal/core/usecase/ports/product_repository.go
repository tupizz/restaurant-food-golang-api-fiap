package ports

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
)

type ProductFilter struct {
	Category string
	Page     int
	PageSize int
}

type ProductRepository interface {
	GetAll(ctx context.Context, filter *ProductFilter) ([]domain.Product, int, error)
	Create(ctx context.Context, product domain.Product) (domain.Product, error)
	GetById(ctx context.Context, id int) (domain.Product, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product domain.Product) (domain.Product, error)
	GetByIds(ctx context.Context, ids []int) ([]domain.Product, int, error)
}
