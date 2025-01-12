package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/shared"
)

type productRepository struct {
	db     *pgxpool.Pool
	sqlcDb *fiapRestaurantDb.Queries
}

func NewProductRepository(
	db *pgxpool.Pool,
	sqlcDb *fiapRestaurantDb.Queries,
) domain.ProductRepository {
	return &productRepository{
		db:     db,
		sqlcDb: sqlcDb,
	}
}

func (r *productRepository) GetByIds(ctx context.Context, ids []int) ([]entity.Product, int, error) {
	if len(ids) == 0 {
		return nil, 0, nil
	}

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
        WHERE p.deleted_at IS NULL AND p.id = ANY($1)
    `

	rows, err := r.db.Query(ctx, query, ids)
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

	return products, len(products), nil
}

func (r *productRepository) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	slog.Info("Updating product", "product", shared.ToJSON(product))

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Product{}, err
	}

	// In case of a panic or error, ensure the transaction is rolled back
	defer func() {
		if p := recover(); p != nil {
			slog.Error("Rolling back transaction", "error", p)
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			slog.Error("Rolling back transaction", "error", err)
			tx.Rollback(ctx)
		} else {
			slog.Info("Committing transaction")
			err = tx.Commit(ctx)
		}
	}()

	// build a dynamic query to update the product
	var columns []string
	var args []any
	argIndex := 1

	if product.Category.Handle != "" {
		// search for category by handle and add id to product
		query := `SELECT id FROM categories WHERE handle = $1`
		err = tx.QueryRow(ctx, query, product.Category.Handle).Scan(&product.Category.ID)
		if err != nil {
			slog.Error("Error searching for category", "error", err)
			return entity.Product{}, domain.ErrNotFound("category")
		}

		// if category not found, throw error
		if product.Category.ID == 0 {
			slog.Error("Category not found", "category", product.Category)
			return entity.Product{}, domain.ErrNotFound("category")
		}

		columns = append(columns, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, product.Category.ID)
		argIndex++
	}

	if product.Name != "" {
		columns = append(columns, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, product.Name)
		argIndex++
	}

	if product.Description != "" {
		columns = append(columns, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, product.Description)
		argIndex++
	}

	if product.Price != 0 {
		columns = append(columns, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, product.Price)
		argIndex++
	}

	if len(columns) > 0 {
		query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", strings.Join(columns, ", "), argIndex)
		slog.Info("Updating product", "query", query)
		args = append(args, product.ID)
		slog.Info("Updating product", "args", args)
		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			slog.Error("Error updating product", "error", err)
			return entity.Product{}, err
		}
	}

	if len(product.Images) > 0 {
		_, err = tx.Exec(ctx, "UPDATE products_images SET deleted_at = NOW() WHERE product_id = $1", product.ID)
		if err != nil {
			slog.Error("Error deleting product images", "error", err)
			return entity.Product{}, err
		}

		for _, image := range product.Images {
			_, err = tx.Exec(ctx, "INSERT INTO products_images (product_id, image) VALUES ($1, $2)", product.ID, image.ImageURL)
			if err != nil {
				slog.Error("Error inserting product image", "error", err)
				return entity.Product{}, err
			}
		}
	}

	product, err = getOneProductWithExecutor(ctx, tx, product.ID)
	if err != nil {
		return entity.Product{}, err
	}

	slog.Info("Updated product", "product", shared.ToJSON(product))

	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	slog.Info("Deleting product", "id", id)
	_, err := r.db.Exec(ctx, "UPDATE products SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	slog.Info("Creating product", "product", shared.ToJSON(product))

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

	// search for category by handle and add id to product
	query := `SELECT id FROM categories WHERE handle = $1`
	err = tx.QueryRow(ctx, query, product.Category.Handle).Scan(&product.Category.ID)
	if err != nil {
		return entity.Product{}, domain.ErrNotFound("category")
	}

	// if category not found, throw error
	if product.Category.ID == 0 {
		return entity.Product{}, domain.ErrNotFound("category")
	}

	query = `INSERT INTO products (name, description, price, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
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
			c.handle AS category_handle,
			c.created_at AS category_created_at,
			c.updated_at AS category_updated_at,
            pi.id AS image_id, 
            pi.image,
			pi.created_at AS image_created_at,
			pi.updated_at AS image_updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL
		LEFT JOIN products_images pi ON p.id = pi.product_id AND pi.deleted_at IS NULL
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
			&product.Category.Handle,
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
