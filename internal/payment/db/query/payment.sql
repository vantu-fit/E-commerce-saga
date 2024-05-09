-- name: GetPaymentById :one
SELECT * FROM payments WHERE id = $1;

-- name: CreatePayment :one
INSERT INTO payments (id, customer_id, currency, amount)
VALUES ($1, $2, $3, $4)
RETURNING *; 

-- name: DeletePayment :one
DELETE FROM payments WHERE id = $1
RETURNING *;