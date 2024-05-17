package command

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
)

type UpdateComment struct {
	*pb.UpdateCommentRequest
	UserID uuid.UUID
}

type UpdateCommentHandler CommandHanlder[UpdateComment]

type updateCommentHandler struct {
	store db.Store
}

func NewUpdateCommentHandler(store db.Store) UpdateCommentHandler {
	return &updateCommentHandler{
		store: store,
	}
}

func (h *updateCommentHandler) Handle(ctx context.Context, cmd UpdateComment) error {
	// check comment
	comment , err := h.store.GetCommentByID(ctx, uuid.MustParse(cmd.Id))
	if err != nil {
		return err
	}
	// check user
	if comment.UserID != cmd.UserID {
		return errors.New("user not allow to update this comment")
	}

	// update comment
	_, err = h.store.UpdateContentComment(ctx, db.UpdateContentCommentParams{
		Content: cmd.Content,
		ID:      uuid.MustParse(cmd.Id),
	})
	return err
}
