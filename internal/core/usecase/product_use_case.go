package usecase

import (
	"context"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type ProductUseCase interface {
	GetProducts(ctx context.Context, filter *ports.ProductFilter) ([]dto.ProductOutput, int, error)
}

type productUseCase struct {
	productRepository ports.ProductRepository
}

func NewProductUseCase(productRepository ports.ProductRepository) ProductUseCase {
	return &productUseCase{productRepository: productRepository}
}

func (c *productUseCase) GetProducts(ctx context.Context, filter *ports.ProductFilter) ([]dto.ProductOutput, int, error) {
	products, total, err := c.productRepository.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	productOutputs := make([]dto.ProductOutput, 0, len(products))
	for _, product := range products {
		if product.ID == 0 {
			continue
		}

		images := make([]string, 0, len(product.Images))
		for _, image := range product.Images {
			if image.ImageURL != "" {
				slog.Info("Processing image URL", "url", image.ImageURL)
				images = append(images, image.ImageURL)
			}
		}

		productOutput := dto.ProductOutput{
			ID:          product.ID,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
			Category:    product.Category.Name,
			Images:      images,
		}

		productOutputs = append(productOutputs, productOutput)
	}

	return productOutputs, total, nil
}
