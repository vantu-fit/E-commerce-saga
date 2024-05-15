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
	*pb.UploadRequest
	Contentype string
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
	image_id := uuid.New()
	err := h.media.UploadObject(ctx, &media.File{
		Data:      bytes.NewReader(cmd.Data),
		Bucket:    media.ImageBucket,
		ProductID: uuid.MustParse(cmd.GetProductId()),
		ID: image_id,
		Contentype: cmd.Contentype,
	})
	if err != nil {
		return "", err
	}
	_, err = h.store.CreateProductImage(ctx, db.CreateProductImageParams{
		ContentType: cmd.Contentype,
		ID:        image_id,
		ProductID: uuid.MustParse(cmd.GetProductId()),
		Alt:       cmd.Alt,
	})
	if err != nil {
		return "", err
	}

	url := "http://" + h.media.GetConfig().Minio.Endpoint + "/" + media.ImageBucket + "/" + image_id.String() + cmd.Contentype

	return url, nil
}
