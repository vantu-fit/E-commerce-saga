package service

import (
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/comment/service/command"
	"github.com/vantu-fit/saga-pattern/internal/comment/service/query"
)

type Command struct {
	CreateComment command.CreateCommentHandler
	DeleteComment command.DeleteCommentHandler
	UpdateComment command.UpdateCommentHandler
}

type Query struct {
	ListComment query.GetCommentsHander
	ListCommentByProductID query.GetCommentsByProductIDHander
}

type Service struct {
	Command Command
	Query   Query
}

func NewService(store db.Store) Service {
	return Service{
		Command: Command{
			CreateComment: command.NewCreateCommentHandler(store),
			DeleteComment: command.NewDeleteCommentHandler(store),
			UpdateComment: command.NewUpdateCommentHandler(store),
		},
		Query: Query{
			ListComment: query.NewGetCommentsHander(store),
			ListCommentByProductID: query.NewGetCommentsByProductHander(store),
		},
	}
}
