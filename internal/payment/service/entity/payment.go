package entity

import "github.com/google/uuid"

type Payment struct {
	ID uuid.UUID
	CurrentcyCode string
	Amount uint64
}