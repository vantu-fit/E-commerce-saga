-- name: CreateOrder :one
INSERT INTO orders (
    id,
    product_id,
    quantity,
    customer_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: DeleteOrder :many
DELETE FROM orders WHERE id = $1 RETURNING *;

-- name: GetOrder :many
SELECT * FROM orders WHERE id = $1;
