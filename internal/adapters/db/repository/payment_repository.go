package repository

import (
	"context"
	"log/slog"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlcDB "github.com/tupizz/restaurant-food-golang-api-fiap/database/sqlc"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/domain"
	"github.com/tupizz/restaurant-food-golang-api-fiap/internal/core/usecase/ports"
)

type paymentRepository struct {
	sqlcDb *sqlcDB.Queries
	dbPool *pgxpool.Pool
}

func NewPaymentRepository(sqlcDb *sqlcDB.Queries, dbPool *pgxpool.Pool) ports.PaymentRepository {
	return &paymentRepository{
		sqlcDb: sqlcDb,
		dbPool: dbPool,
	}
}

func (r *paymentRepository) UpdateOrderPaymentStatus(ctx context.Context, externalReference string, paymentMethod string, status domain.PaymentStatus) error {
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
			slog.Error("Order update failed", "externalRefrence", externalReference, "paymentMethod", paymentMethod, "status", status)
			slog.Error("Transaction", "tx", tx)
			slog.Error("Context", "ctx", ctx)
			slog.Error("Error", "error", err)
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	qtx := r.sqlcDb.WithTx(tx)

	err = qtx.UpdateOrderPaymentStatus(ctx, sqlcDB.UpdateOrderPaymentStatusParams{
		ExternalReference: pgtype.Text{
			String: externalReference,
			Valid:  true,
		},
		Method: paymentMethod,
		Status: pgtype.Text{
			String: string(status),
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	var orderStatusToUpdate domain.OrderStatus
	if status == domain.PaymentStatusApproved {
		orderStatusToUpdate = domain.OrderStatusPreparing
	} else {
		orderStatusToUpdate = domain.OrderStatusCanceled
	}

	orderId, err := qtx.GetOrderIdByExternalReferenceAndMethod(ctx, sqlcDB.GetOrderIdByExternalReferenceAndMethodParams{
		ExternalReference: pgtype.Text{
			String: externalReference,
			Valid:  true,
		},
		Method: paymentMethod,
	})
	if err != nil {
		return err
	}

	err = qtx.UpdateOrderStatus(ctx, sqlcDB.UpdateOrderStatusParams{
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
