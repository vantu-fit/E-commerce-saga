package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/pb"
)

type GetVideoUrlsHandler QueryHandler[*pb.GetUrlRequest, *pb.GetUrlResponse]

type getVideoUrlsHandler struct {
	config *config.Config
	store  db.Store
}

func NewGetVideoUrlsHandler(
	config *config.Config,
	store db.Store,
) GetVideoUrlsHandler {
	return &getVideoUrlsHandler{
		config: config,
		store:  store,
	}
}

func (h *getVideoUrlsHandler) Handle(ctx context.Context, req *pb.GetUrlRequest) (*pb.GetUrlResponse, error) {
	videos, err := h.store.GetProductVideoByProductID(ctx, uuid.MustParse(req.GetProductId()))
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(videos))
	for i, image := range videos {
		urls[i] = "http://" + h.config.Minio.Endpoint + "/" + media.VideoBucket + "/" + image.ID.String() + image.ContentType
	}

	response := pb.GetUrlResponse{
		Url: urls,
	}

	return &response, nil
}
