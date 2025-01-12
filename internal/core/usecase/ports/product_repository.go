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
}
