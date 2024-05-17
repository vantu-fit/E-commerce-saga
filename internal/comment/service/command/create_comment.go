package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
)

type CreateComment struct {
	*pb.CreateCommentRequest
	UserID uuid.UUID
}

type CreateCommentHandler CommandHanlder[CreateComment]

type createCommentHandler struct {
	store db.Store
}

func NewCreateCommentHandler(store db.Store) CreateCommentHandler {
	return &createCommentHandler{
		store: store,
	}
}

func (h *createCommentHandler) Handle(ctx context.Context, cmd CreateComment) error {
	err := h.store.CreateCommentTx(ctx, db.CreateCommentParamsTx{
		CreateCommentRequest: cmd.CreateCommentRequest,
		UserID:               cmd.UserID,
	})
	return err
}
