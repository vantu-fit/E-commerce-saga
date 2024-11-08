// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Comment struct {
	ID         uuid.UUID   `json:"id"`
	ProductID  uuid.UUID   `json:"product_id"`
	UserID     uuid.UUID   `json:"user_id"`
	Content    string      `json:"content"`
	LeftIndex  int64       `json:"left_index"`
	RightIndex int64       `json:"right_index"`
	ParentID   pgtype.UUID `json:"parent_id"`
	UpadatedAt time.Time   `json:"upadated_at"`
	CreatedAt  time.Time   `json:"created_at"`
}
