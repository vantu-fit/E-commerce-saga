package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/product/service/entity"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type UpdateProductInventory struct {
	PurchaseID uuid.UUID
	ProductItems *[]entity.ProductItem
}

type UpdateProductInventoryHandler CommandHanlder[UpdateProductInventory]

type updateProductInventoryHandler struct {
	store db.Store
}

func NewUpdateProductInventoryHandler(store db.Store) UpdateProductInventoryHandler {
	return &updateProductInventoryHandler{
		store: store,
	}
}

// idempentency key is purchase id

func (h *updateProductInventoryHandler) Handle (ctx context.Context , cmd UpdateProductInventory) error {
	purchaseItems := make([]db.PurchasedProduct , len(*cmd.ProductItems))
	for i , product := range *cmd.ProductItems {
		purchaseItems[i] = db.PurchasedProduct{
			ProductID: product.ID,
			Quantity: product.Quantity,
		}
	}
	return h.store.UpdateProductInventoryTx(ctx , cmd.PurchaseID , &purchaseItems)
}