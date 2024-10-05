package repository_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/di"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

func TestProductRepository_Create(t *testing.T) {
	// Build the container
	container := di.BuildContainer()

	// Invoke the test function
	err := container.Invoke(func(productRepository domain.ProductRepository) {
		// Prepare the product to create
		productToCreate := entity.Product{
			Name:        "Produto novo",
			Description: "Descrição do produto novo",
			Price:       10.0,
			Category: entity.ProductCategory{
				ID: 1,
			},
			Images: []entity.ProductImage{
				{
					ImageURL: "https://placehold.co/600x400/png",
				},
			},
		}

		// Create the product
		createdProduct, err := productRepository.Create(context.Background(), productToCreate)
		assert.NoError(t, err, "Error creating product")

		// Assert that the returned product has the expected values
		assert.NotZero(t, createdProduct.ID, "Product ID should not be zero")
		assert.Equal(t, productToCreate.Name, createdProduct.Name, "Product names should match")
		assert.Equal(t, productToCreate.Description, createdProduct.Description, "Product descriptions should match")
		assert.Equal(t, productToCreate.Price, createdProduct.Price, "Product prices should match")
		assert.Equal(t, productToCreate.Category.ID, createdProduct.Category.ID, "Category IDs should match")
		assert.NotEmpty(t, createdProduct.Images, "Product should have images")
		assert.Equal(t, productToCreate.Images[0].ImageURL, createdProduct.Images[0].ImageURL, "Image URLs should match")

		// Optionally, print the product as pretty JSON for debugging
		productJSON, err := json.MarshalIndent(createdProduct, "", "  ")
		assert.NoError(t, err, "Error marshalling product to JSON")

		fmt.Println(string(productJSON))
	})

	// Handle errors from the container invocation
	if err != nil {
		t.Fatalf("Error invoking container: %v", err)
	}
}

func TestProductRepository_GetById(t *testing.T) {
	// Build the container
	container := di.BuildContainer()

	// Invoke the test function
	err := container.Invoke(func(productRepository domain.ProductRepository) {
		// Create a product first
		productToCreate := entity.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       15.0,
			Category: entity.ProductCategory{
				ID: 1,
			},
			Images: []entity.ProductImage{
				{
					ImageURL: "https://example.com/test.jpg",
				},
			},
		}

		createdProduct, err := productRepository.Create(context.Background(), productToCreate)
		assert.NoError(t, err, "Error creating product")

		// Get the product by ID
		retrievedProduct, err := productRepository.GetById(context.Background(), createdProduct.ID)
		assert.NoError(t, err, "Error getting product by ID")

		// Assert that the retrieved product matches the created product
		assert.Equal(t, createdProduct.ID, retrievedProduct.ID, "Product IDs should match")
		assert.Equal(t, createdProduct.Name, retrievedProduct.Name, "Product names should match")
		assert.Equal(t, createdProduct.Description, retrievedProduct.Description, "Product descriptions should match")
		assert.Equal(t, createdProduct.Price, retrievedProduct.Price, "Product prices should match")
		assert.Equal(t, createdProduct.Category.ID, retrievedProduct.Category.ID, "Category IDs should match")
		assert.Equal(t, len(createdProduct.Images), len(retrievedProduct.Images), "Number of images should match")
		assert.Equal(t, createdProduct.Images[0].ImageURL, retrievedProduct.Images[0].ImageURL, "Image URLs should match")

		// Optionally, print the retrieved product as pretty JSON for debugging
		productJSON, err := json.MarshalIndent(retrievedProduct, "", "  ")
		assert.NoError(t, err, "Error marshalling product to JSON")
		fmt.Println(string(productJSON))
	})

	// Handle errors from the container invocation
	if err != nil {
		t.Fatalf("Error invoking container: %v", err)
	}
}
func TestProductRepository_GetAll(t *testing.T) {
	// Build the container
	container := di.BuildContainer()

	// Invoke the test function
	err := container.Invoke(func(productRepository domain.ProductRepository) {
		// Create multiple products
		for i := 0; i < 15; i++ {
			productToCreate := entity.Product{
				Name:        fmt.Sprintf("Test Product %d", i+1),
				Description: fmt.Sprintf("Test Description %d", i+1),
				Price:       float64(10 + i),
				Category: entity.ProductCategory{
					ID: 1,
				},
				Images: []entity.ProductImage{
					{
						ImageURL: fmt.Sprintf("https://example.com/test%d.jpg", i+1),
					},
				},
			}

			_, err := productRepository.Create(context.Background(), productToCreate)
			assert.NoError(t, err, "Error creating product")
		}

		// Test pagination
		page := 1
		pageSize := 10

		products, totalCount, err := productRepository.GetAll(context.Background(), page, pageSize)
		assert.NoError(t, err, "Error getting all products")

		// Assert pagination results
		assert.LessOrEqual(t, len(products), pageSize, "Number of products should not exceed page size")
		assert.Greater(t, totalCount, 0, "Total count should be greater than 0")

		// Assert product details
		for _, product := range products {
			fmt.Println(product)
			assert.NotZero(t, product.ID, "Product ID should not be zero")
			assert.NotEmpty(t, product.Name, "Product name should not be empty")
			assert.NotEmpty(t, product.Description, "Product description should not be empty")
			assert.Greater(t, product.Price, 0.0, "Product price should be greater than 0")
			assert.NotZero(t, product.Category.ID, "Category ID should not be zero")
		}

		// Optionally, print the products as pretty JSON for debugging
		productsJSON, err := json.MarshalIndent(products, "", "  ")
		assert.NoError(t, err, "Error marshalling products to JSON")
		fmt.Println(string(productsJSON))
		fmt.Printf("Total count: %d\n", totalCount)
	})

	// Handle errors from the container invocation
	if err != nil {
		t.Fatalf("Error invoking container: %v", err)
	}
}
