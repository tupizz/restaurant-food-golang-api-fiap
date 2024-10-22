-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;


-- name: UpdateOrderPaymentStatus :exec
UPDATE payments
SET status = $2
WHERE order_id = $1;
