package grpc

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// check req
	violations := val.NewValidator("category_id", req.GetCategoryId()).UUID().Validate()
	violations = append(violations, val.NewValidator("name", req.GetName()).String().MinLenght(3).MaxLenght(255).Validate()...)
	violations = append(violations, val.NewValidator("description", req.GetDescription()).String().MinLenght(3).MaxLenght(255).Validate()...)
	violations = append(violations, val.NewValidator("brand_name", req.GetBrandName()).String().MinLenght(3).MaxLenght(255).Validate()...)
	// violations = append(violations, val.NewValidator("price", req.GetPrice()).Number().Min(0).Validate()...)
	// violations = append(violations, val.NewValidator("inventory", req.GetInventory()).Number().Min(0).Validate()...)
	if violations != nil {
		return nil, val.InvalidArgumentError(violations)
	}

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
	authResponse, err := server.grpcClient.AccountClient.Auth(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	if !authResponse.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	// create category
	arg := db.CreateProductParams{
		IDCategory:  uuid.MustParse(req.GetCategoryId()),
		IDAccount:   uuid.MustParse(authResponse.UserId),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		BrandName:   req.GetBrandName(),
		Price:       req.GetPrice(),
		Inventory:   req.GetInventory(),
	}

	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create category: %s", err)
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          product.ID.String(),
			CategoryId:  product.IDCategory.String(),
			IdAccount:   product.IDAccount.String(),
			Name:        product.Name,
			Description: product.Description,
			BrandName:   product.BrandName,
			Price:       product.Price,
			Inventory:   product.Inventory,
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
			CreatedAt:   timestamppb.New(product.CreatedAt),
		},
	}, nil

}
