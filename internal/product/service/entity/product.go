package entity

import "github.com/google/uuid"

type ProductItem struct {
	ID       uuid.UUID
	Quantity int64
}
