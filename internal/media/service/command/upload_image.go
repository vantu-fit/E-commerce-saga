package command

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/pb"
)

type UploadImage struct {
	*pb.UploadImageRequest
}

type UploadImageHandler CommandHandler[UploadImage, string]

type uploadImageHandler struct {
	store db.Store
	media media.Media
}

func NewUploadImageHandler(
	store db.Store,
	media media.Media,
) UploadImageHandler {
	return &uploadImageHandler{
		store: store,
		media: media,
	}
}

func (h *uploadImageHandler) Handle(ctx context.Context, cmd UploadImage) (string, error) {
	err := h.media.UploadObject(ctx, &media.File{
		Name:      cmd.Filename,
		Data:      bytes.NewReader(cmd.Data),
		Bucket:    media.ImageBucket,
		ProductID: uuid.MustParse(cmd.GetProductId()),
	})
	if err != nil {
		return "", err
	}

	_, err = h.store.CreateProductImage(ctx, db.CreateProductImageParams{
		ProductID: uuid.MustParse(cmd.GetProductId()),
		Name:      cmd.Filename,
		Alt:       cmd.Alt,
	})
	if err != nil {
		return "", err
	}

	url := "http://" + h.media.GetConfig().Minio.Endpoint + "/" + media.ImageBucket + "/" + cmd.Filename

	return url, nil
}
