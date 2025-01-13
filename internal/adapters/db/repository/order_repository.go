package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	sqlcDB "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type orderRepository struct {
	db     *pgxpool.Pool
	sqlcDb *sqlcDB.Queries
}

func NewOrderRepository(db *pgxpool.Pool, sqlcDb *sqlcDB.Queries) ports.OrderRepository {
	return &orderRepository{
		db:     db,
		sqlcDb: sqlcDb,
	}
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func (r *orderRepository) GetAll(ctx context.Context, filter *ports.OrderFilter) ([]sqlcDB.GetAllOrdersRow, error) {
	orders, err := r.sqlcDb.GetAllOrders(ctx, sqlcDB.GetAllOrdersParams{
		Limit:  int32(filter.PageSize),
		Offset: int32((filter.Page - 1) * filter.PageSize),
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer tx.Rollback(ctx)

	// Create Order
	query := `
		INSERT INTO orders (client_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, order.ClientID, order.Status, time.Now(), time.Now()).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return domain.Order{}, err
	}

	// Create Order Items
	for idx, item := range order.Items {
		item.OrderID = order.ID
		createdItem, err := r.createOrderItem(ctx, tx, &item)
		if err != nil {
			return domain.Order{}, err
		}
		order.Items[idx] = *createdItem
	}

	// Create Payment
	order.Payment.OrderID = order.ID
	_, err = r.createPayment(ctx, tx, &order.Payment)
	if err != nil {
		return domain.Order{}, err
	}

	// Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (domain.Order, error) {
	// Fetch Order
	query := `
		SELECT id, client_id, status, created_at, updated_at, deleted_at
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`
	var order domain.Order
	err := r.db.QueryRow(ctx, query, id).
		Scan(&order.ID, &order.ClientID, &order.Status, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt)
	if err == pgx.ErrNoRows {
		return domain.Order{}, ErrOrderNotFound
	} else if err != nil {
		return domain.Order{}, err
	}

	// Fetch Order Items
	order.Items, err = r.getOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return domain.Order{}, err
	}

	// Fetch Payment
	order.Payment, err = r.getPaymentByOrderID(ctx, order.ID)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) Update(ctx context.Context, order domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update Order
	query := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	_, err = tx.Exec(ctx, query, order.Status, time.Now(), order.ID)
	if err != nil {
		return err
	}

	// Update Order Items
	for _, item := range order.Items {
		_, err := r.updateOrderItem(ctx, tx, item)
		if err != nil {
			return err
		}
	}

	// Update Payment
	err = r.updatePayment(ctx, tx, order.Payment)
	if err != nil {
		return err
	}

	// Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Soft delete the order
	query := `
		UPDATE orders
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err = tx.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	// Soft delete order items
	err = r.deleteOrderItemsByOrderID(ctx, tx, id)
	if err != nil {
		return err
	}

	// Soft delete payment
	err = r.deletePaymentByOrderID(ctx, tx, id)
	if err != nil {
		return err
	}

	// Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) createOrderItem(ctx context.Context, tx pgx.Tx, item *domain.OrderItem) (*domain.OrderItem, error) {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.Price, time.Now(), time.Now()).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *orderRepository) updateOrderItem(ctx context.Context, tx pgx.Tx, item domain.OrderItem) (domain.OrderItem, error) {
	query := `
		UPDATE order_items
		SET quantity = $1, price = $2, updated_at = $3
		WHERE id = $4 AND order_id = $5
		RETURNING id, updated_at
	`
	err := tx.QueryRow(ctx, query, item.Quantity, item.Price, time.Now(), item.ID, item.OrderID).
		Scan(&item.ID, &item.UpdatedAt)
	if err != nil {
		return domain.OrderItem{}, err
	}
	return item, nil
}

func (r *orderRepository) deleteOrderItemsByOrderID(ctx context.Context, tx pgx.Tx, orderID int) error {
	query := `
		UPDATE order_items
		SET deleted_at = $1
		WHERE order_id = $2 AND deleted_at IS NULL
	`
	_, err := tx.Exec(ctx, query, time.Now(), orderID)
	return err
}

func (r *orderRepository) getOrderItemsByOrderID(ctx context.Context, orderID int) ([]domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price, created_at, updated_at, deleted_at
		FROM order_items
		WHERE order_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *orderRepository) createPayment(ctx context.Context, tx pgx.Tx, payment *domain.Payment) (*domain.Payment, error) {
	query := `
		INSERT INTO payments (order_id, status, method, amount, external_reference, qr_data ,created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(ctx, query, payment.OrderID, payment.Status, payment.Method, payment.Amount, payment.ExternalReference, payment.QRData, time.Now(), time.Now()).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *orderRepository) updatePayment(ctx context.Context, tx pgx.Tx, payment domain.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, method = $2, amount = $3, updated_at = $4
		WHERE id = $5 AND order_id = $6
	`
	_, err := tx.Exec(ctx, query, payment.Status, payment.Method, payment.Amount, time.Now(), payment.ID, payment.OrderID)
	return err
}

func (r *orderRepository) deletePaymentByOrderID(ctx context.Context, tx pgx.Tx, orderID int) error {
	query := `
		UPDATE payments
		SET deleted_at = $1
		WHERE order_id = $2 AND deleted_at IS NULL
	`
	_, err := tx.Exec(ctx, query, time.Now(), orderID)
	return err
}

func (r *orderRepository) getPaymentByOrderID(ctx context.Context, orderID int) (domain.Payment, error) {
	query := `
		SELECT id, order_id, status, method, amount, created_at, updated_at, deleted_at
		FROM payments
		WHERE order_id = $1 AND deleted_at IS NULL
	`
	var payment domain.Payment
	err := r.db.QueryRow(ctx, query, orderID).
		Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.Method, &payment.Amount, &payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt)
	if err == pgx.ErrNoRows {
		return domain.Payment{}, ErrOrderNotFound
	} else if err != nil {
		return domain.Payment{}, err
	}

	return payment, nil
}
