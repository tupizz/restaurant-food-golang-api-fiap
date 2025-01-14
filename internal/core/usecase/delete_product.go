package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type DeleteProductUseCase interface {
	Run(ctx context.Context, id int) error
}

type deleteProductUseCase struct {
	productRepository ports.ProductRepository
}

func NewDeleteProductUseCase(productRepository ports.ProductRepository) DeleteProductUseCase {
	return &deleteProductUseCase{productRepository: productRepository}
}

func (c *deleteProductUseCase) Run(ctx context.Context, id int) error {
	return c.productRepository.Delete(ctx, id)
}
