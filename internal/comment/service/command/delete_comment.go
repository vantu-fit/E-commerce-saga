package command

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
)

type DeleteComment struct {
	*pb.DeleteCommentRequest
	UserID uuid.UUID
}

type DeleteCommentHandler CommandHanlder[DeleteComment]

type deleteCommentHandler struct {
	store db.Store
}

func NewDeleteCommentHandler(store db.Store) DeleteCommentHandler {
	return &deleteCommentHandler{
		store: store,
	}
}

func (h *deleteCommentHandler) Handle(ctx context.Context, cmd DeleteComment) error {
	// check comment
	comment, err := h.store.GetCommentByID(ctx, uuid.MustParse(cmd.Id))
	if err != nil {
		return err
	}

	// check user
	if comment.UserID != cmd.UserID {
		return errors.New("user not allow to delete this comment")
	}

	// delete comment
	err = h.store.DeleteCommentTx(ctx, db.DeleteCommentParamsTx{
		DeleteCommentRequest: cmd.DeleteCommentRequest,
	})
	return err
}
