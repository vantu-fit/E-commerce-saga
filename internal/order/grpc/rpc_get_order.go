package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/order/service/query"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetOrderById(ctx context.Context, req *pb.GetOrderByIdRequest) (*pb.GetOrderByIdResponse, error) {
	// check uuid
	orderId, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// query order detail
	orderDetail, err := server.service.GetOrderDetail.Handle(ctx, query.GetOrderDetail{
		OrderID: orderId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// make product to response
	detailProduct := make([]*pb.DetailOrderProduct, len(*orderDetail.PurchasedProducts))

	response := &pb.GetOrderByIdResponse{
		Id:         orderId.String(),
		CustomerId: orderDetail.CustomerID.String(),
		Products:   detailProduct,
		CreatedAt:  timestamppb.Now(), // TODO: return time create order
	}

	for i, product := range *orderDetail.PurchasedProducts {
		detailProduct[i] = &pb.DetailOrderProduct{
			Id:          product.ID.String(),
			CategoryId:  product.CategoryID.String(),
			Name:        product.Name,
			BrandName:   product.BrandName,
			Description: product.Description,
			Price:       uint64(product.Price),
			Quantity:    uint64(product.Quantity),
		}
	}

	return response, nil
}
