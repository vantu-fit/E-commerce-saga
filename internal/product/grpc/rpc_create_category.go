package grpc

import (
	"context"

	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	// check req
	violations := val.NewValidator("name", req.GetName()).String().MinLenght(3).MaxLenght(50).Validate()
	violations = append(violations, val.NewValidator("description", req.GetDescription()).String().MinLenght(3).MaxLenght(255).Validate()...)
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
	authResponse, err := server.grpcClient.accountClient.Auth(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	if !authResponse.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	// create category
	arg := db.CreateCategoryParams{
		Name:        req.Name,
		Description: req.Description,
	}

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create category: %s", err)
	}

	return &pb.CreateCategoryResponse{
		Categoty: &pb.Category{
			Id:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
			UpdatedAt:   timestamppb.New(category.UpdatedAt),
			CreatedAt:   timestamppb.New(category.CreatedAt),
		},
	}, nil

}
