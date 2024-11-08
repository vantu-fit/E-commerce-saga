// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: video.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createProductVideo = `-- name: CreateProductVideo :one
INSERT INTO product_videos (
    id,
    content_type,
    product_id,
    alt
) VALUES (
    $1, $2, $3 , $4
) RETURNING id, content_type, product_id, alt, created_at
`

type CreateProductVideoParams struct {
	ID          uuid.UUID `json:"id"`
	ContentType string    `json:"content_type"`
	ProductID   uuid.UUID `json:"product_id"`
	Alt         string    `json:"alt"`
}

func (q *Queries) CreateProductVideo(ctx context.Context, arg CreateProductVideoParams) (ProductVideo, error) {
	row := q.db.QueryRow(ctx, createProductVideo,
		arg.ID,
		arg.ContentType,
		arg.ProductID,
		arg.Alt,
	)
	var i ProductVideo
	err := row.Scan(
		&i.ID,
		&i.ContentType,
		&i.ProductID,
		&i.Alt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProductVideoByID = `-- name: DeleteProductVideoByID :one
DELETE FROM product_videos WHERE id = $1 RETURNING id, content_type, product_id, alt, created_at
`

func (q *Queries) DeleteProductVideoByID(ctx context.Context, id uuid.UUID) (ProductVideo, error) {
	row := q.db.QueryRow(ctx, deleteProductVideoByID, id)
	var i ProductVideo
	err := row.Scan(
		&i.ID,
		&i.ContentType,
		&i.ProductID,
		&i.Alt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProductVideoByProductID = `-- name: DeleteProductVideoByProductID :many
DELETE FROM product_videos WHERE product_id = $1 RETURNING id, content_type, product_id, alt, created_at
`

func (q *Queries) DeleteProductVideoByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVideo, error) {
	rows, err := q.db.Query(ctx, deleteProductVideoByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductVideo
	for rows.Next() {
		var i ProductVideo
		if err := rows.Scan(
			&i.ID,
			&i.ContentType,
			&i.ProductID,
			&i.Alt,
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

const getProductVideoByID = `-- name: GetProductVideoByID :one
SELECT id, content_type, product_id, alt, created_at FROM product_videos WHERE id = $1
`

func (q *Queries) GetProductVideoByID(ctx context.Context, id uuid.UUID) (ProductVideo, error) {
	row := q.db.QueryRow(ctx, getProductVideoByID, id)
	var i ProductVideo
	err := row.Scan(
		&i.ID,
		&i.ContentType,
		&i.ProductID,
		&i.Alt,
		&i.CreatedAt,
	)
	return i, err
}

const getProductVideoByProductID = `-- name: GetProductVideoByProductID :many
SELECT id, content_type, product_id, alt, created_at FROM product_videos WHERE product_id = $1
`

func (q *Queries) GetProductVideoByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVideo, error) {
	rows, err := q.db.Query(ctx, getProductVideoByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductVideo
	for rows.Next() {
		var i ProductVideo
		if err := rows.Scan(
			&i.ID,
			&i.ContentType,
			&i.ProductID,
			&i.Alt,
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
