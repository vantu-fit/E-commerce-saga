package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/internal/media/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	// Auth

	// Validate request

	// Call command
	url, err := s.service.Commad.UploadImage.Handle(ctx, command.UploadImage{
		UploadImageRequest: req,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Response
	respose := &pb.UploadImageResponse{
		Url: url,
	}
	return respose, nil
}
