package service

import (
	"context"
	"log/slog"

	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/application/dto"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/shared"
)

type ProductServiceAdmin interface {
	CreateProduct(ctx context.Context, input dto.ProductInputCreate) (dto.ProductOutput, error)
	UpdateProduct(ctx context.Context, id int, input dto.ProductInputUpdate) (dto.ProductOutput, error)
	DeleteProduct(ctx context.Context, id int) error
}

type productServiceAdmin struct {
	productRepository domain.ProductRepository
}

func NewProductServiceAdmin(productRepo domain.ProductRepository) ProductServiceAdmin {
	return &productServiceAdmin{productRepository: productRepo}
}

func (s *productServiceAdmin) DeleteProduct(ctx context.Context, id int) error {
	return s.productRepository.Delete(ctx, id)
}

func (s *productServiceAdmin) UpdateProduct(ctx context.Context, id int, input dto.ProductInputUpdate) (dto.ProductOutput, error) {
	images := make([]entity.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, entity.ProductImage{
			ImageURL: image,
		})
	}

	product := entity.Product{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: entity.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	updatedProduct, err := s.productRepository.Update(ctx, product)
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

func (s *productServiceAdmin) CreateProduct(ctx context.Context, input dto.ProductInputCreate) (dto.ProductOutput, error) {
	images := make([]entity.ProductImage, 0, len(input.Images))
	for _, image := range input.Images {
		images = append(images, entity.ProductImage{
			ImageURL: image,
		})
	}

	product := entity.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category: entity.ProductCategory{
			Handle: input.Category,
		},
		Images: images,
	}

	createdProduct, err := s.productRepository.Create(ctx, product)
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
