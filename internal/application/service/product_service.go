package service

import (
	"context"
	"log"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
)

type ProductService interface {
	GetProducts(ctx context.Context, page int, pageSize int) ([]dto.ProductOutput, int, error)
}

type productService struct {
	productRepository domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) ProductService {
	return &productService{productRepository: productRepo}
}

func (s *productService) GetProducts(ctx context.Context, page int, pageSize int) ([]dto.ProductOutput, int, error) {
	products, total, err := s.productRepository.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Create an empty slice, but with the capacity of the number of products
	// otherwise we can face items filled with zero values
	productOutputs := make([]dto.ProductOutput, 0, len(products))
	for _, product := range products {
		if product.ID == 0 {
			continue
		}

		// Create an empty slice, but with the capacity of the number of images
		// otherwise we can face items filled with zero values
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

		log.Println(productOutput)

		slog.Info("Created product output", "product", productOutput)

		productOutputs = append(productOutputs, productOutput)
	}

	return productOutputs, total, nil
}
