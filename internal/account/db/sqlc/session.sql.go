// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: session.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  id ,
  user_id,
  refresh_token,
  user_agent,
  client_ip
) VALUES (
    $1 , $2, $3, $4, $5 
) RETURNING id, user_id, refresh_token, user_agent, client_ip
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
	)
	return i, err
}

const getSessionById = `-- name: GetSessionById :one
SELECT id, user_id, refresh_token, user_agent, client_ip FROM sessions WHERE id = $1
`

func (q *Queries) GetSessionById(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSessionById, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
	)
	return i, err
}
