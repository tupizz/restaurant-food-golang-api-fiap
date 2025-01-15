package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type UpdateProductUseCase interface {
	Run(ctx context.Context, id int, input dto.ProductInputUpdate) (*entities.Product, error)
}

type updateProductUseCase struct {
	productRepository ports.ProductRepository
}

func NewUpdateProductUseCase(productRepository ports.ProductRepository) UpdateProductUseCase {
	return &updateProductUseCase{productRepository: productRepository}
}

func (c *updateProductUseCase) Run(ctx context.Context, id int, input dto.ProductInputUpdate) (*entities.Product, error) {
	images := make([]entities.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, entities.ProductImage{
			ImageURL: image,
		})
	}

	product := entities.Product{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: entities.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	updatedProduct, err := c.productRepository.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}
