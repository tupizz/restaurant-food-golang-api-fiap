-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: UpdateOrderPaymentStatus :exec
UPDATE payments
SET status = $3
WHERE external_reference = $1 AND method = $2;

-- name: GetOrderIdByExternalReferenceAndMethod :one
SELECT order_id
FROM payments
WHERE external_reference = $1 AND method = $2;
