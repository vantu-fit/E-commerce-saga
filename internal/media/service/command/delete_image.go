package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/pb"
)

type DeleteImage struct {
	*pb.DeleteImageRequest
}

type DeleteImageHandler CommandHandler[DeleteImage, any]

type deleteImageHandler struct {
	store db.Store
	media media.Media
}

func NewDeleteImageHandler(
	store db.Store,
	media media.Media,
) DeleteImageHandler {
	return &deleteImageHandler{
		store: store,
		media: media,
	}
}

func (h *deleteImageHandler) Handle(ctx context.Context, cmd DeleteImage) (any, error) {
	fileInfo, err := h.media.GetOjectInfo(ctx, &media.File{
		Name:      cmd.Filename,
		Bucket:    media.ImageBucket,
		ProductID: uuid.New(),
		Data:      nil,
	})
	if err != nil {
		return nil, err
	}

	err = h.media.DeleteObject(ctx, &media.File{})
	if err != nil {
		return nil, err
	}

	product_id := fileInfo.UserMetadata["product_id"]

	_, err = h.store.DeleteProductImageByID(ctx, uuid.MustParse(product_id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
