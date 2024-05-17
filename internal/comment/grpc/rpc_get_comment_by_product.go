package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/vantu-fit/saga-pattern/internal/comment/service/query"
	"github.com/vantu-fit/saga-pattern/pb"
)

func (s *Server) GetCommentByProduct(ctx context.Context, req *pb.GetCommentByProductRequest) (*pb.GetCommentByProductResponse, error) {
	res, err := s.service.Query.ListCommentByProductID.Handle(ctx , query.GetCommentsByProductID{
		GetCommentByProductRequest: req,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error, err: %v", err)
	}

	return res.GetCommentByProductResponse, nil
}
