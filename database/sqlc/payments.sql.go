// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: payments.sql

package fiapRestaurantDb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const updateOrderPaymentStatus = `-- name: UpdateOrderPaymentStatus :exec
UPDATE payments
SET status = $2
WHERE order_id = $1
`

type UpdateOrderPaymentStatusParams struct {
	OrderID int32
	Status  pgtype.Text
}

func (q *Queries) UpdateOrderPaymentStatus(ctx context.Context, arg UpdateOrderPaymentStatusParams) error {
	_, err := q.db.Exec(ctx, updateOrderPaymentStatus, arg.OrderID, arg.Status)
	return err
}

const updateOrderStatus = `-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1
`

type UpdateOrderStatusParams struct {
	ID     int32
	Status pgtype.Text
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error {
	_, err := q.db.Exec(ctx, updateOrderStatus, arg.ID, arg.Status)
	return err
}
