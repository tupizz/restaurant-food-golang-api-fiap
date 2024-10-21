-- name: GetAllOrders :many
WITH paginated_orders AS (
    SELECT DISTINCT ON (o.id) o.id, o.created_at
    FROM orders o
    WHERE o.deleted_at IS NULL
    LIMIT $1
    OFFSET $2
)
select o.id         as order_id,
       o.created_at as order_created_at,
       o.updated_at as order_updated_at,
       o.status     as order_status,
       c.name       as client_name,
       c.cpf as client_cpf,
       c.id         as client_id,
       p.id         as product_id,
       p.name       as product_name,
       p.price      as product_price,
       p.description as product_description,
       oi.quantity  as product_quantity,
       py.id as payment_id,
       py.status as payment_status,
       py.amount as payment_amount,
       py.method as payment_method,
       pt.handle    as category_handle
from paginated_orders po 
         join orders o on o.id = po.id
         join order_items oi on oi.order_id = o.id
         join products p on oi.product_id = p.id
         join categories pt on p.category_id = pt.id
         join payments py on py.order_id = o.id
         join clients c on c.id = o.client_id
WHERE o.deleted_at IS NULL
ORDER BY o.created_at DESC;
