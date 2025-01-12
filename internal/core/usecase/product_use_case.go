package usecase

import (
	"context"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/shared"
)

type ProductUseCase interface {
	GetProducts(ctx context.Context, filter *ports.ProductFilter) ([]dto.ProductOutput, int, error)
	CreateProduct(ctx context.Context, input dto.ProductInputCreate) (dto.ProductOutput, error)
	UpdateProduct(ctx context.Context, id int, input dto.ProductInputUpdate) (dto.ProductOutput, error)
	DeleteProduct(ctx context.Context, id int) error
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

func (c *productUseCase) DeleteProduct(ctx context.Context, id int) error {
	return c.productRepository.Delete(ctx, id)
}

func (c *productUseCase) UpdateProduct(ctx context.Context, id int, input dto.ProductInputUpdate) (dto.ProductOutput, error) {
	images := make([]domain.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, domain.ProductImage{
			ImageURL: image,
		})
	}

	product := domain.Product{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: domain.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	updatedProduct, err := c.productRepository.Update(ctx, product)
	if err != nil {
		slog.Error("Error updating product", "error", err)
		return dto.ProductOutput{}, err
	}

	slog.Info("Updated product", "product", shared.ToJSON(updatedProduct))

	imageOutputs := make([]string, 0, len(updatedProduct.Images))
	for _, image := range updatedProduct.Images {
		imageOutputs = append(imageOutputs, image.ImageURL)
	}

	return dto.ProductOutput{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		Category:    updatedProduct.Category.Handle,
		Images:      imageOutputs,
	}, nil
}

func (c *productUseCase) CreateProduct(ctx context.Context, input dto.ProductInputCreate) (dto.ProductOutput, error) {
	images := make([]domain.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, domain.ProductImage{
			ImageURL: image,
		})
	}

	product := domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: domain.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	createdProduct, err := c.productRepository.Create(ctx, product)
	if err != nil {
		return dto.ProductOutput{}, err
	}

	slog.Info("Created product", "product", shared.ToJSON(createdProduct))

	imageOutputs := make([]string, 0, len(createdProduct.Images))
	for _, image := range createdProduct.Images {
		imageOutputs = append(imageOutputs, image.ImageURL)
	}

	return dto.ProductOutput{
		ID:          createdProduct.ID,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		Category:    createdProduct.Category.Handle,
		Images:      imageOutputs,
	}, nil
}
