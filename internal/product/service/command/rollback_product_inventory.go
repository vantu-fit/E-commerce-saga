package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/product/service/entity"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type RollbackProductInventory struct {
	PurchaseID uuid.UUID
	ProductItems *[]entity.ProductItem
}

type RollbackProductInventoryHandler CommandHanlder[RollbackProductInventory]

type rollbackProductInventoryHandler struct {
	store db.Store
}

func NewRollbackProductInventoryHandler(store db.Store) RollbackProductInventoryHandler {
	return &rollbackProductInventoryHandler{
		store: store,
	}
}

// idempentency key is purchase id

func (h *rollbackProductInventoryHandler) Handle (ctx context.Context , cmd RollbackProductInventory) error {
	purchaseItems := make([]db.PurchasedProduct , len(*cmd.ProductItems))
	for i , product := range *cmd.ProductItems {
		purchaseItems[i] = db.PurchasedProduct{
			ProductID: product.ID,
			Quantity: product.Quantity,
		}
	}
	return h.store.RollbackProductInventoryTx(ctx , cmd.PurchaseID , &purchaseItems)
}