package query

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
)

type CheckProduct struct {
	*pb.CheckProductRequest
}

type CheckProductResult struct {
	*pb.CheckProductResponse
}

type CheckProductHandler QueryHandler[CheckProduct, CheckProductResult]

type checkProductHandler struct {
	store db.Store
}

func NewCheckProductHandler(store db.Store) CheckProductHandler {
	return &checkProductHandler{
		store: store,
	}
}

func (h *checkProductHandler) Handle(ctx context.Context, cmd CheckProduct) (CheckProductResult, error) {
	res, err := h.store.GetProductInventory(ctx, uuid.MustParse(cmd.Id))
	if err != nil {
		return CheckProductResult{
			CheckProductResponse: &pb.CheckProductResponse{
				Valid: false,
			},
		}, err
	}
	if res.Inventory < cmd.Quantity {
		log.Info().Msgf("Product: product %d is out of stock , Quantity %d", res.Inventory , cmd.Quantity)
		return CheckProductResult{
			CheckProductResponse: &pb.CheckProductResponse{
				Valid: false,
			},
		}, nil
	}

	return CheckProductResult{
		CheckProductResponse: &pb.CheckProductResponse{
			Valid: true,
		},
	}, nil

}
