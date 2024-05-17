package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/internal/comment/service/query"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	res, err := s.service.Query.ListComment.Handle(ctx, query.GetComments{
		GetCommentRequest: req,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error, err: %v", err)
	}

	return res.GetCommentResponse, nil
}
