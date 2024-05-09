package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
)


type UpdateProductDetailHandler CommandHanlder[*pb.UpdateProductRequest]

type updateProductDetailHandler struct {
	store db.Store
}

func NewUpdateProductDetailHandler(store db.Store) UpdateProductDetailHandler {
	return &updateProductDetailHandler{
		store: store,
	}
}

func (h *updateProductDetailHandler) Handle(ctx context.Context, cmd *pb.UpdateProductRequest) error {
	arg := db.UpadateProductParams{
		ID: uuid.MustParse(cmd.GetId()),
		IDCategory: pgtype.UUID{
			Bytes: uuid.MustParse(cmd.GetCategoryId()),
			Valid: cmd.CategoryId != nil,
		},
		Name: pgtype.Text{
			String: cmd.GetName(),
			Valid:  cmd.Name != nil,
		},
		Description: pgtype.Text{
			String: cmd.GetDescription(),
			Valid:  cmd.Description != nil,
		},
		BrandName: pgtype.Text{
			String: cmd.GetBrandName(),
			Valid:  cmd.BrandName != nil,
		},
		Price: pgtype.Int8{
			Int64: cmd.GetPrice(),
			Valid: cmd.Price != nil,
		},
		Inventory: pgtype.Int8{
			Int64: cmd.GetInventory(),
			Valid: cmd.Inventory != nil,
		},
	}
	_, err := h.store.UpadateProduct(ctx , arg )
	return err
}
