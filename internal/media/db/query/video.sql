-- name: CreateProductVideo :one
INSERT INTO product_videos (
    id,
    content_type,
    product_id,
    alt
) VALUES (
    $1, $2, $3 , $4
) RETURNING *;
-- name: DeleteProductVideoByID :one
DELETE FROM product_videos WHERE id = $1 RETURNING *;
-- name: DeleteProductVideoByProductID :many
DELETE FROM product_videos WHERE product_id = $1 RETURNING *;
-- name: GetProductVideoByID :one
SELECT * FROM product_videos WHERE id = $1;
-- name: GetProductVideoByProductID :many
SELECT * FROM product_videos WHERE product_id = $1;
