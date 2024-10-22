package repository

import (
	"context"
	"log/slog"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	fiapRestaurantDb "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/domain/entity"
)

type paymentRepository struct {
	sqlcDb *fiapRestaurantDb.Queries
	dbPool *pgxpool.Pool
}

func NewPaymentRepository(
	sqlcDb *fiapRestaurantDb.Queries,
	dbPool *pgxpool.Pool,
) domain.PaymentRepository {
	return &paymentRepository{
		sqlcDb: sqlcDb,
		dbPool: dbPool,
	}
}

func (r *paymentRepository) UpdateOrderPaymentStatus(ctx context.Context, orderId int, status entity.PaymentStatus) error {
	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Error starting transaction", "error", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			slog.Error("Order update failed", "order_id", orderId, "status", status)
			slog.Error("Transaction", "tx", tx)
			slog.Error("Context", "ctx", ctx)
			slog.Error("Error", "error", err)
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	qtx := r.sqlcDb.WithTx(tx)

	err = qtx.UpdateOrderPaymentStatus(ctx, fiapRestaurantDb.UpdateOrderPaymentStatusParams{
		OrderID: int32(orderId),
		Status: pgtype.Text{
			String: string(status),
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	var orderStatusToUpdate entity.OrderStatus
	if status == entity.PaymentStatusApproved {
		orderStatusToUpdate = entity.OrderStatusPreparing
	} else {
		orderStatusToUpdate = entity.OrderStatusCanceled
	}

	err = qtx.UpdateOrderStatus(ctx, fiapRestaurantDb.UpdateOrderStatusParams{
		ID: int32(orderId),
		Status: pgtype.Text{
			String: string(orderStatusToUpdate),
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
