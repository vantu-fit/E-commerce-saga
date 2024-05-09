-- name: CreateProductVideo :one
INSERT INTO product_videos (
    product_id,
    name,
    alt
) VALUES (
    $1, $2, $3
) RETURNING *;
-- name: DeleteProductVideoByID :one
DELETE FROM product_videos WHERE id = $1 RETURNING *;
-- name: DeleteProductVideoByProductID :many
DELETE FROM product_videos WHERE product_id = $1 RETURNING *;
-- name: GetProductVideoByID :one
SELECT * FROM product_videos WHERE id = $1;
-- name: GetProductVideoByProductID :many
SELECT * FROM product_videos WHERE product_id = $1;
