package object


import "github.com/google/uuid"

type PurchasedProduct struct {
	ID uuid.UUID
	Quantity int64
}

type DetailedPurchasedProduct struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	BrandName   string
	Description string
	Price       uint32
	Quantity    uint32
}