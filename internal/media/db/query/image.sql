-- name: CreateProductImage :one
INSERT INTO product_images (
    id,
    content_type,
    product_id,
    alt
) VALUES (
    $1, $2, $3 , $4
) RETURNING *;
-- name: DeleteProductImageByID :one
DELETE FROM product_images WHERE id = $1 RETURNING *;
-- name: DeleteProductImageByProductID :many
DELETE FROM product_images WHERE product_id = $1 RETURNING *;
-- name: GetProductImageByID :one
SELECT * FROM product_images WHERE id = $1;
-- name: GetProductImagesByProductID :many
SELECT * FROM product_images WHERE product_id = $1;
-- name: UpdateProductImage :one
