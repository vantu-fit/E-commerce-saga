-- name: CreateProduct :one 
INSERT INTO products (
  id_category,
  id_account,
  name,
  description,
  brand_name,
  price,
  inventory
) VALUES (
  $1, $2, $3, $4, $5, $6 , $7
) RETURNING *;

-- name: GetProductCategory :many
SELECT * FROM products JOIN categories ON products.id_category = categories.id WHERE categories.name = $1;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductForUpdate :one
SELECT * FROM products WHERE id = $1 LIMIT 1 FOR UPDATE;

-- name: UpadateProduct :one
UPDATE products
SET
  id_category = COALESCE(sqlc.narg(id_category), id_category),
  name = COALESCE(sqlc.narg(name), name),
  description = COALESCE(sqlc.narg(description), description),
  brand_name = COALESCE(sqlc.narg(brand_name), brand_name),
  price = COALESCE(sqlc.narg(price), price),
  inventory = COALESCE(sqlc.narg(inventory), inventory),
  updated_at = now()
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: UpdateProductInventory :one
UPDATE products
SET
  inventory = inventory + $2
WHERE
  id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id OFFSET $1 LIMIT $2;

-- name: GetProductInventory :one
SELECT id,inventory FROM products WHERE id = $1;