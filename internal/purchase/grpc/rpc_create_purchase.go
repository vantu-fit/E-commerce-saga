package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreatePurchase(ctx context.Context, req *pb.CreatePurchaseRequestApi) (*pb.CreatePurchaseResponseApi, error) {
	fmt.Println("CreatePurchase")
	//Authclient
	// transform context to metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	value := md.Get("Authorization")
	if len(value) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing token")
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", value[0])
	// end transform context to metadata
	authResponse, err := s.grpcClient.AccountClient.Auth(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	if !authResponse.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}
	// ProductClient
	// count amount
	amount := int64(0)
	for _, item := range req.OrderItems {
		checkRes ,err := s.grpcClient.ProductClient.CheckProduct(ctx, &pb.CheckProductRequest{
			Id: item.ProductId,
			Quantity: int64(item.Quantity),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Error check product, err: %s", err.Error())
		}

		if !checkRes.Valid {
			return nil, status.Errorf(codes.InvalidArgument, "Product out of stock")
		}

		product, err := s.grpcClient.ProductClient.GetProductByID(ctx, &pb.GetProductByIDRequest{
			Id: item.ProductId,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Error get product, err: %s", err.Error())
		}

		amount += product.Product.Price * int64(item.Quantity)
	
	}

	purchaseID := uuid.New()
	err = s.service.Command.CreatePurchase.Handle(ctx, command.CreatePurchase{
		CreatePurchaseRequestApi: req,
		CustomerId:               uuid.MustParse(authResponse.UserId),
		Amount:                   amount,
		PurchaseId:               purchaseID,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error create purchase, err: %s", err.Error())
	}

	return &pb.CreatePurchaseResponseApi{
		PurchaseId: purchaseID.String(),
		Status:     "pending",
		Timestamp:  timestamppb.Now(),
	}, nil
}
