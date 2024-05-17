package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/pb"
)

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "pong",
	}, nil
}
