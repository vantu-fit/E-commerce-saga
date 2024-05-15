package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/internal/product/service/query"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CheckProduct(ctx context.Context, req *pb.CheckProductRequest) (*pb.CheckProductResponse, error) {
	res, err := server.service.Query.CheckProduct.Handle(ctx, query.CheckProduct{
		CheckProductRequest: req,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error check product, err: %s", err.Error())
	}
	return &pb.CheckProductResponse{
		Valid: res.Valid,
	}, nil
}
