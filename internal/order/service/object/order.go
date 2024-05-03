package object

import "github.com/google/uuid"

type DetailedOrder struct {
	ID                uuid.UUID
	CustomerID        uuid.UUID
	PurchasedProducts *[]DetailedPurchasedProduct
}
