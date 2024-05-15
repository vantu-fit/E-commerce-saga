package command

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/pb"
)

type UploadVideo struct {
	*pb.UploadRequest
	Contentype string
}

type UploadVideoHandler CommandHandler[UploadVideo, string]

type uploadVideoHandler struct {
	store db.Store
	media media.Media
}

func NewUploadVideoHandler(
	store db.Store,
	media media.Media,
) UploadVideoHandler {
	return &uploadVideoHandler{
		store: store,
		media: media,
	}
}

func (h *uploadVideoHandler) Handle(ctx context.Context, cmd UploadVideo) (string, error) {
	image_id := uuid.New()
	err := h.media.UploadObject(ctx, &media.File{
		Data:      bytes.NewReader(cmd.Data),
		Bucket:    media.VideoBucket,
		ProductID: uuid.MustParse(cmd.GetProductId()),
		ID: image_id,
		Contentype: cmd.Contentype,
	})
	if err != nil {
		return "", err
	}
	_, err = h.store.CreateProductVideo(ctx, db.CreateProductVideoParams{
		ContentType: cmd.Contentype,
		ID:        image_id,
		ProductID: uuid.MustParse(cmd.GetProductId()),
		Alt:       cmd.Alt,
	})
	if err != nil {
		return "", err
	}

	url := "http://" + h.media.GetConfig().Minio.Endpoint + "/" + media.VideoBucket + "/" + image_id.String() + cmd.Contentype

	return url, nil
}
