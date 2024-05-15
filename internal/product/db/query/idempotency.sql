-- name: GetIdempotencyKey :many
SELECT * FROM idempotency WHERE id = $1 ; 

-- name: CreateIdempotency :one
INSERT INTO idempotency (
  id,
  product_id,
  quantity,
  rollbacked
) VALUES (
  $1, $2, $3 , $4
) RETURNING *;

-- name: UpdateIdempotency :many
UPDATE idempotency
SET
  rollbacked = COALESCE(sqlc.narg(rollbacked), rollbacked)
WHERE
  id = sqlc.arg(id) 
RETURNING *;
