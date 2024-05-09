package entity

import "github.com/google/uuid"

type Purchase struct {
	ID      uuid.UUID
	Order   *Order
	Payment *Payment
}

