// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: idempotency.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createIdempotency = `-- name: CreateIdempotency :one
INSERT INTO idempotency (
  id,
  product_id,
  quantity,
  rollbacked
) VALUES (
  $1, $2, $3 , $4
) RETURNING id, product_id, quantity, rollbacked, created_at
`

type CreateIdempotencyParams struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int32     `json:"quantity"`
	Rollbacked bool      `json:"rollbacked"`
}

func (q *Queries) CreateIdempotency(ctx context.Context, arg CreateIdempotencyParams) (Idempotency, error) {
	row := q.db.QueryRow(ctx, createIdempotency,
		arg.ID,
		arg.ProductID,
		arg.Quantity,
		arg.Rollbacked,
	)
	var i Idempotency
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Quantity,
		&i.Rollbacked,
		&i.CreatedAt,
	)
	return i, err
}

const getIdempotencyKey = `-- name: GetIdempotencyKey :many
SELECT id, product_id, quantity, rollbacked, created_at FROM idempotency WHERE id = $1
`

func (q *Queries) GetIdempotencyKey(ctx context.Context, id uuid.UUID) ([]Idempotency, error) {
	rows, err := q.db.Query(ctx, getIdempotencyKey, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Idempotency
	for rows.Next() {
		var i Idempotency
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Quantity,
			&i.Rollbacked,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateIdempotency = `-- name: UpdateIdempotency :many
UPDATE idempotency
SET
  rollbacked = COALESCE($1, rollbacked)
WHERE
  id = $2 
RETURNING id, product_id, quantity, rollbacked, created_at
`

type UpdateIdempotencyParams struct {
	Rollbacked pgtype.Bool `json:"rollbacked"`
	ID         uuid.UUID   `json:"id"`
}

func (q *Queries) UpdateIdempotency(ctx context.Context, arg UpdateIdempotencyParams) ([]Idempotency, error) {
	rows, err := q.db.Query(ctx, updateIdempotency, arg.Rollbacked, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Idempotency
	for rows.Next() {
		var i Idempotency
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Quantity,
			&i.Rollbacked,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}