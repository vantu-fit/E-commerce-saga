package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetVideoUrl(ctx context.Context, req *pb.GetUrlRequest) (*pb.GetUrlResponse, error) {
	// AUTH
	response, err := s.service.Query.GetVideoUrls.Handle(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}
