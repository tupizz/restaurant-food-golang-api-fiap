package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	// Begin a new transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Product{}, err
	}

	// In case of a panic or error, ensure the transaction is rolled back
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Rolling back transaction")
			tx.Rollback(ctx)
			panic(p) // Re-throw the panic after rolling back
		} else if err != nil {
			fmt.Println("Rolling back transaction")
			tx.Rollback(ctx)
		} else {
			fmt.Println("Committing transaction")
			err = tx.Commit(ctx)
		}
	}()

	query := `INSERT INTO products (name, description, price, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.Category.ID).Scan(&product.ID)
	if err != nil {
		return entity.Product{}, err
	}

	query = `INSERT INTO products_images (product_id, image) VALUES ($1, $2)`
	for _, image := range product.Images {
		_, err = tx.Exec(ctx, query, product.ID, image.ImageURL)
		if err != nil {
			return entity.Product{}, err
		}
	}

	// get product with images
	query = `
		SELECT 
			p.id, 
			p.name, 
			p.description, 
			p.price, 
			c.id AS category_id, 
			c.name AS category_name, 
			pi.id AS image_id, 
			pi.image
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN products_images pi ON p.id = pi.product_id
		WHERE p.id = $1
	`
	rows, err := tx.Query(ctx, query, product.ID)
	if err != nil {
		return product, err
	}
	defer rows.Close()

	var result_product entity.Product
	for rows.Next() {
		var (
			id           int
			name         string
			description  string
			price        float64
			categoryID   int
			categoryName string
			imageID      sql.NullInt64
			imageURL     sql.NullString
		)

		err = rows.Scan(&id, &name, &description, &price, &categoryID, &categoryName, &imageID, &imageURL)
		if err != nil {
			return product, err
		}

		if result_product.ID == 0 {
			result_product = entity.Product{
				ID:          id,
				Name:        name,
				Description: description,
				Price:       price,
				Category: entity.ProductCategory{
					ID:   categoryID,
					Name: categoryName,
				},
				Images: []entity.ProductImage{},
			}
		}

		// Append images if present
		if imageID.Valid && imageURL.Valid {
			result_product.Images = append(result_product.Images, entity.ProductImage{
				ID:       int(imageID.Int64),
				ImageURL: imageURL.String,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return product, err
	}

	return result_product, nil
}

func (r *productRepository) GetById(ctx context.Context, id int) (entity.Product, error) {
	return entity.Product{}, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	return []entity.Product{}, nil
}
