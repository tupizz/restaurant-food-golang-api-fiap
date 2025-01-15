-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, category_id = $5
WHERE id = $1
RETURNING *;

-- name: GetProductsByIds :many
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
WHERE p.deleted_at IS NULL AND p.id = ANY(@ids::int[]);

-- name: GetProductById :many
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
WHERE p.id = $1 AND p.deleted_at IS NULL;
