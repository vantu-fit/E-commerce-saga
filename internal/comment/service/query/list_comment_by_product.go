package query

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetCommentsByProductID struct {
	*pb.GetCommentByProductRequest
}

type GetCommentsByProductIDResult struct {
	*pb.GetCommentByProductResponse
}

type GetCommentsByProductIDHander QueryHandler[GetCommentsByProductID, *GetCommentsByProductIDResult]

type getCommentsByProductIDHander struct {
	store db.Store
}

func NewGetCommentsByProductHander(store db.Store) GetCommentsByProductIDHander {
	return &getCommentsByProductIDHander{
		store: store,
	}
}

func (h *getCommentsByProductIDHander) Handle(ctx context.Context, cmd GetCommentsByProductID) (*GetCommentsByProductIDResult, error) {
	var result = GetCommentsByProductIDResult{}

	// get comments by product id
	comments, err := h.store.GetCommentByProductID(ctx, uuid.MustParse(cmd.ProductId))
	if err != nil {
		return nil, err
	}

	commentsRes := make([]*pb.Comment, len(comments))
	for i, c := range comments {
		commentsRes[i] = &pb.Comment{
			Id:        c.ID.String(),
			ProductId: c.ProductID.String(),
			Content:   c.Content,
			ParentId:  uuid.UUID(c.ParentID.Bytes).String(),
			UpdatedAt: timestamppb.New(c.UpadatedAt),
			CreatedAt: timestamppb.New(c.CreatedAt),
		}
	}

	result.GetCommentByProductResponse = &pb.GetCommentByProductResponse{
		Comment: commentsRes,
	}

	return &result, nil
}
