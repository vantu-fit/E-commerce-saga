package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/pb"
)

type GetImageUrlsHandler QueryHandler[*pb.GetUrlRequest, *pb.GetUrlResponse]

type getImageUrlsHandler struct {
	config *config.Config
	store  db.Store
}

func NewGetImageUrlsHandler(
	config *config.Config,
	store db.Store,
) GetImageUrlsHandler {
	return &getImageUrlsHandler{
		config: config,
		store:  store,
	}
}

func (h *getImageUrlsHandler) Handle(ctx context.Context, req *pb.GetUrlRequest) (*pb.GetUrlResponse, error) {
	images, err := h.store.GetProductImagesByProductID(ctx, uuid.MustParse(req.GetProductId()))
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(images))
	for i, image := range images {
		urls[i] = "http://" + h.config.Minio.Endpoint + "/" + media.ImageBucket + "/" + image.ID.String() + image.ContentType
	}

	response := pb.GetUrlResponse{
		Url: urls,
	}

	return &response, nil
}
