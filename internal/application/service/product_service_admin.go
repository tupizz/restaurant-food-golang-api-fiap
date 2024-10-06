package service

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
)

type ProductServiceAdmin interface {
	CreateProduct(ctx context.Context, input dto.ProductInput) (dto.ProductOutput, error)
	UpdateProduct(ctx context.Context, id int, input dto.ProductInput) (dto.ProductOutput, error)
	DeleteProduct(ctx context.Context, id int) error
}
