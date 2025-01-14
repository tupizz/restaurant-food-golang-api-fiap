package usecase

import (
	"context"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type CreateProductUseCase interface {
	Run(ctx context.Context, input dto.ProductInputCreate) (*entities.Product, error)
}

type createProductsUseCase struct {
	productRepository ports.ProductRepository
}

func NewCreateProductUseCase(productRepository ports.ProductRepository) CreateProductUseCase {
	return &createProductsUseCase{productRepository: productRepository}
}

func (c *createProductsUseCase) Run(ctx context.Context, input dto.ProductInputCreate) (*entities.Product, error) {
	images := make([]entities.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, entities.ProductImage{
			ImageURL: image,
		})
	}

	product := entities.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: entities.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	createdProduct, err := c.productRepository.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &createdProduct, nil
}
