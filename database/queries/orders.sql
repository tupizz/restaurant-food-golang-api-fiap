-- name: GetAllOrders :many
WITH paginated_orders AS (
    SELECT DISTINCT ON (o.id) o.id, o.created_at
    FROM orders o
    WHERE o.deleted_at IS NULL 
      AND o.status IN ('ready', 'preparing', 'pending')
    ORDER BY o.id, o.created_at DESC
    LIMIT $1
    OFFSET $2
)
SELECT o.id         AS order_id,
       o.created_at AS order_created_at,
       o.updated_at AS order_updated_at,
       o.status     AS order_status,
       c.name       AS client_name,
       c.cpf        AS client_cpf,
       c.id         AS client_id,
       p.id         AS product_id,
       p.name       AS product_name,
       p.price      AS product_price,
       p.description AS product_description,
       oi.quantity  AS product_quantity,
       py.id        AS payment_id,
       py.status    AS payment_status,
       py.amount    AS payment_amount,
       py.method    AS payment_method,
       pt.handle    AS category_handle
FROM paginated_orders po 
JOIN orders o ON o.id = po.id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON oi.product_id = p.id
JOIN categories pt ON p.category_id = pt.id
JOIN payments py ON py.order_id = o.id
JOIN clients c ON c.id = o.client_id
WHERE o.deleted_at IS NULL AND o.status IN ('ready', 'preparing', 'pending')
ORDER BY 
    CASE 
        WHEN o.status = 'ready' THEN 1
        WHEN o.status = 'preparing' THEN 2
        WHEN o.status = 'pending' THEN 3
        ELSE 4
    END,
    o.created_at DESC;
