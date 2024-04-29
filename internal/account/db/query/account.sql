-- name: CreateAccount :one
INSERT INTO accounts (
  first_name,
  last_name,
  email,
  address,
  phone_number,
  password
) VALUES (
  $1, $2, $3 , $4, $5, $6
) RETURNING *;


-- name: GetAccountByEmail :one
SELECT * FROM accounts WHERE email = $1;

