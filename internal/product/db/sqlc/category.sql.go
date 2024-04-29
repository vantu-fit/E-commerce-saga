// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: category.sql

package db

import (
	"context"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories ( 
  name,
  description
) VALUES (
  $1, $2 
) RETURNING id, name, description, updated_at, created_at
`

type CreateCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCategory, arg.Name, arg.Description)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
