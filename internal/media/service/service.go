package service

import (
	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/internal/media/service/command"
	"github.com/vantu-fit/saga-pattern/internal/media/service/query"
)

type Commad struct {
	UploadImage command.UploadImageHandler
	DeleteImage command.DeleteImageHandler
	UploadVideo command.UploadVideoHandler
}
type Query struct {
	GetImageUrls query.GetImageUrlsHandler
	GetVideoUrls query.GetVideoUrlsHandler
}

type Service struct {
	Commad
	Query
}

func NewService(
	config *config.Config,
	store db.Store,
	media media.Media,
) *Service {
	return &Service{
		Commad: Commad{
			UploadImage: command.NewUploadImageHandler(store, media),
			DeleteImage: command.NewDeleteImageHandler(store, media),
			UploadVideo: command.NewUploadVideoHandler(store, media),
		},
		Query: Query{
			GetImageUrls: query.NewGetImageUrlsHandler(config, store),
			GetVideoUrls: query.NewGetVideoUrlsHandler(config, store),
		},
	}
}
