package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
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

func (r *productRepository) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r.db.Exec(ctx, "UPDATE products SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5", product.Name, product.Description, product.Price, product.Category.ID, product.ID)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
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

	product, err = getOneProductWithExecutor(ctx, tx, product.ID)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (r *productRepository) GetById(ctx context.Context, id int) (entity.Product, error) {
	product, err := getOneProductWithExecutor(ctx, r.db, id)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, page int, pageSize int) ([]entity.Product, int, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query to get paginated results
	query := `
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
        ORDER BY p.id
    `
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	productMap := make(map[int]*entity.Product)
	var products []entity.Product

	for rows.Next() {
		var product entity.Product
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
			// Product already exists, just add the image
			if imageID.Valid && imageURL.Valid {
				existingProduct.Images = append(existingProduct.Images, entity.ProductImage{
					ID:        int(imageID.Int64),
					ImageURL:  imageURL.String,
					CreatedAt: imageCreatedAt.Time,
					UpdatedAt: imageUpdatedAt.Time,
				})
			}
		} else {
			// New product, add it to the map and slice
			if imageID.Valid && imageURL.Valid {
				product.Images = []entity.ProductImage{{
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

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(DISTINCT id) FROM products`
	err = r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func getOneProductWithExecutor(ctx context.Context, executor interface{}, id int) (entity.Product, error) {
	query := `
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
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN products_images pi ON p.id = pi.product_id
		WHERE p.id = $1 AND p.deleted_at IS NULL
	`
	var rows pgx.Rows
	var err error

	switch e := executor.(type) {
	case pgx.Tx:
		rows, err = e.Query(ctx, query, id)
	case *pgxpool.Pool:
		rows, err = e.Query(ctx, query, id)
	default:
		return entity.Product{}, fmt.Errorf("unsupported executor type")
	}

	if err != nil {
		return entity.Product{}, err
	}

	defer rows.Close()

	var result_product entity.Product
	for rows.Next() {
		var product entity.Product
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
			return entity.Product{}, err
		}

		if result_product.ID == 0 {
			result_product = product
		}

		if imageID.Valid && imageURL.Valid {
			result_product.Images = append(result_product.Images, entity.ProductImage{
				ID:        int(imageID.Int64),
				ImageURL:  imageURL.String,
				CreatedAt: imageCreatedAt.Time,
				UpdatedAt: imageUpdatedAt.Time,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return entity.Product{}, err
	}

	return result_product, nil
}
