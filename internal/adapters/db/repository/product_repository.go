package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ports.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(ctx context.Context, filter *ports.ProductFilter) ([]domain.Product, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize

	categoryFilter := ""
	if filter.Category != "" {
		categoryFilter = fmt.Sprintf("AND LOWER(c.handle) = LOWER('%s')", filter.Category)
	}

	query := fmt.Sprintf(`
			WITH paginated_products AS (
				SELECT DISTINCT ON (p.id) p.id
				FROM products p
				WHERE p.deleted_at IS NULL
				ORDER BY p.id
				LIMIT $1 OFFSET $2
			)
			SELECT 
				p.id, 
				p.name, 
				p.description, 
				p.price, 
				p.created_at,
				p.updated_at,
				c.id AS category_id, 
				c.name AS category_name,
				c.created_at AS category_created_at,
				c.updated_at AS category_updated_at,
				pi.id AS image_id, 
				pi.image,
				pi.created_at AS image_created_at,
				pi.updated_at AS image_updated_at
			FROM paginated_products pp
			JOIN products p ON pp.id = p.id
			LEFT JOIN categories c ON p.category_id = c.id
			LEFT JOIN products_images pi ON p.id = pi.product_id
			WHERE p.deleted_at IS NULL
			%s
			ORDER BY p.id
		`, categoryFilter)

	rows, err := r.db.Query(ctx, query, filter.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	productMap := make(map[int]*domain.Product)
	var products []domain.Product

	for rows.Next() {
		var product domain.Product
		var imageID sql.NullInt64
		var imageURL sql.NullString
		var imageCreatedAt sql.NullTime
		var imageUpdatedAt sql.NullTime

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category.ID,
			&product.Category.Name,
			&product.Category.CreatedAt,
			&product.Category.UpdatedAt,
			&imageID,
			&imageURL,
			&imageCreatedAt,
			&imageUpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		if existingProduct, ok := productMap[product.ID]; ok {
			if imageID.Valid && imageURL.Valid {
				existingProduct.Images = append(existingProduct.Images, domain.ProductImage{
					ID:        int(imageID.Int64),
					ImageURL:  imageURL.String,
					CreatedAt: imageCreatedAt.Time,
					UpdatedAt: imageUpdatedAt.Time,
				})
			}
		} else {
			if imageID.Valid && imageURL.Valid {
				product.Images = []domain.ProductImage{{
					ID:        int(imageID.Int64),
					ImageURL:  imageURL.String,
					CreatedAt: imageCreatedAt.Time,
					UpdatedAt: imageUpdatedAt.Time,
				}}
			}
			productMap[product.ID] = &product
			products = append(products, product)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var totalCount int
	countQuery := `SELECT COUNT(DISTINCT id) FROM products`
	err = r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
