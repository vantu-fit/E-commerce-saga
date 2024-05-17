package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vantu-fit/saga-pattern/internal/comment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetComments struct {
	*pb.GetCommentRequest
}

type GetCommentsResult struct {
	*pb.GetCommentResponse
}

type GetCommentsHander QueryHandler[GetComments, *GetCommentsResult]

type getCommentsHander struct {
	store db.Store
}

func NewGetCommentsHander(store db.Store) GetCommentsHander {
	return &getCommentsHander{
		store: store,
	}
}

func (h *getCommentsHander) Handle(ctx context.Context, cmd GetComments) (*GetCommentsResult, error) {
	var result = GetCommentsResult{}

	// get comments by id
	comment, err := h.store.GetCommentByID(ctx, uuid.MustParse(cmd.Id))
	if err != nil {
		return nil, err
	}

	// check is root comment
	if !comment.ParentID.Valid {
		comments, err := h.store.GetAllComments(ctx, db.GetAllCommentsParams{
			ParentID: pgtype.UUID{
				Bytes: comment.ID,
				Valid: true,
			},
			LeftIndex:  comment.LeftIndex,
			RightIndex: comment.RightIndex,
		})
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
				CreatedAt: timestamppb.New(c.CreatedAt),
				UpdatedAt: timestamppb.New(c.UpadatedAt),
			}
		}

		result.GetCommentResponse = &pb.GetCommentResponse{
			Comment: commentsRes,
		}

		return &result, nil

	}

	comments, err := h.store.GetAllComments(ctx, db.GetAllCommentsParams{
		ParentID: comment.ParentID,
		LeftIndex:  comment.LeftIndex,
		RightIndex: comment.RightIndex,
	})
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
			CreatedAt: timestamppb.New(c.CreatedAt),
			UpdatedAt: timestamppb.New(c.UpadatedAt),
		}
	}

	result.GetCommentResponse = &pb.GetCommentResponse{
		Comment: commentsRes,
	}

	return &result, nil
}
