package service

import (
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/internal/media/service/command"
)

type Commad struct {
	UploadImage command.UploadImageHandler
	DeleteImage command.DeleteImageHandler
}
type Query struct {
}

type Service struct {
	Commad
	Query
}

func NewService(
	store db.Store,
	media media.Media,
) *Service {
	return &Service{
		Commad: Commad{
			UploadImage: command.NewUploadImageHandler(store, media),
			DeleteImage: command.NewDeleteImageHandler(store, media),
		},
	}
}
