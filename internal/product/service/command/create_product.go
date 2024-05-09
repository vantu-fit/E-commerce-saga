package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type CreateProduct struct {
	Name        string
	BrandName   string
	Description string
	Price       int64
	CategoryID  uuid.UUID
	Inventory   int64
	IDAccount   uuid.UUID
}

type CreateProductHandler CommandHanlder[CreateProduct]

type createProductHandler struct {
	store db.Store
}

func NewCreateProductHandler(store db.Store) CreateProductHandler {
	return &createProductHandler{
		store: store,
	}
}

func (h *createProductHandler) Handle(ctx context.Context, cmd CreateProduct) error {
	_, err := h.store.CreateProduct(ctx, db.CreateProductParams{
		IDCategory:  cmd.CategoryID,
		Name:        cmd.Name,
		BrandName:   cmd.BrandName,
		Description: cmd.Description,
		Price:       cmd.Price,
		Inventory:   cmd.Inventory,
		IDAccount:   cmd.IDAccount,
	})
	return err
}
