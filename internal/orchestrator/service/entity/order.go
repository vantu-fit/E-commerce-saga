package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID uuid.UUID
	Quantity uint64
}

type Order struct {
	ID uuid.UUID
	CustomerID uuid.UUID
	OrderItems *[]OrderItem
}