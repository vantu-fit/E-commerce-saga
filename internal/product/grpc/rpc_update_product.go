package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	// check req
	violations := val.NewValidator("id", req.GetId()).UUID().Validate()
	if violations != nil {
		violations = append(violations, val.NewValidator("category_id", req.GetCategoryId()).UUID().Validate()...)
	}
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
	arg := db.UpadateProductParams{
		ID: uuid.MustParse(req.GetId()),
		IDCategory: pgtype.UUID{
			Bytes: uuid.MustParse(req.GetCategoryId()),
			Valid: req.CategoryId != nil,
		},
		Name: pgtype.Text{
			String: req.GetName(),
			Valid:  req.Name != nil,
		},
		Description: pgtype.Text{
			String: req.GetDescription(),
			Valid:  req.Description != nil,
		},
		BrandName: pgtype.Text{
			String: req.GetBrandName(),
			Valid:  req.BrandName != nil,
		},
		Price: pgtype.Int4{
			Int32: req.GetPrice(),
			Valid: req.Price != nil,
		},
		Inventory: pgtype.Int4{
			Int32: req.GetInventory(),
			Valid: req.Inventory != nil,
		},
	}

	product, err := server.store.UpadateProduct(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create category: %s", err)
	}

	return &pb.UpdateProductResponse{
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
