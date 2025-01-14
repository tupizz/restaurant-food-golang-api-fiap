package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type GetProductsUseCase interface {
	Run(ctx context.Context, filter *ports.ProductFilter) ([]entities.Product, int, error)
}

type getProductsUseCase struct {
	productRepository ports.ProductRepository
}

func NewGetProductsUseCase(productRepository ports.ProductRepository) GetProductsUseCase {
	return &getProductsUseCase{productRepository: productRepository}
}

func (c *getProductsUseCase) Run(ctx context.Context, filter *ports.ProductFilter) ([]entities.Product, int, error) {
	products, total, err := c.productRepository.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
