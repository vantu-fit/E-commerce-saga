package command

import (
	"context"

	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type CreateCategory struct {
	Name        string
	Description string
}

type CreateCategoryHandler CommandHanlder[CreateCategory]

type createCategoryHandler struct {
	store db.Store
}

func NewCreateCategoryHandler(store db.Store) CreateCategoryHandler {
	return &createCategoryHandler{
		store: store,
	}
}

func (h *createCategoryHandler) Handle(ctx context.Context, cmd CreateCategory) error {
	_, err := h.store.CreateCategory(ctx, db.CreateCategoryParams{
		Name:        cmd.Name,
		Description: cmd.Description,
	})
	return err
}
