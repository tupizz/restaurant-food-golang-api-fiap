package mappers

import (
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain/entities"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/dto"
)

func ToProductsDTO(products []entities.Product) []dto.ProductOutput {
	productOutputs := make([]dto.ProductOutput, 0, len(products))
	for _, product := range products {
		if product.ID == 0 {
			continue
		}

		productOutputs = append(productOutputs, ToProductDTO(product))
	}

	return productOutputs
}

func ToProductDTO(product entities.Product) dto.ProductOutput {
	images := make([]string, 0, len(product.Images))
	for _, image := range product.Images {
		if image.ImageURL != "" {
			images = append(images, image.ImageURL)
		}
	}

	return dto.ProductOutput{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Category:    product.Category.Name,
		Images:      images,
	}
}
