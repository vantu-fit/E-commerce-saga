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
	fmt.Println(req.OrderItems)
	fmt.Println(req.Payment)
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
	amount := int64(100)
	// for _, item := range req.OrderItems {
	// 	product, err := s.grpcClient.ProductClient.GetProductByID(ctx, &pb.GetProductByIDRequest{
	// 		Id: item.ProductId,
	// 	})
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "Error get product, err: %s", err.Error())
	// 	}

	// 	amount += product.Product.Price * int64(item.Quantity)
	// }
	fmt.Println("amount: ", amount)
	purchaseID := uuid.New()
	err = s.service.Command.CreatePurchase.Handle(ctx, command.CreatePurchase{
		CreatePurchaseRequestApi: req,
		CustomerId:               uuid.MustParse(authResponse.Account.Id),
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
