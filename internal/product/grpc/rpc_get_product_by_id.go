package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/pb"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	// check req
	violations := val.NewValidator("id", req.GetId()).UUID().Validate()
	if violations != nil {
		return nil, val.InvalidArgumentError(violations)
	}

	product, err := server.store.GetProductByID(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "cannot get product: %s", err)
	}
	// get image
	image , err := server.grpcClient.MediaClient.GetImageUrl(ctx , &pb.GetUrlRequest{
		ProductId: product.ID.String(),
	})
	if err != nil {
		log.Error().Msgf("cannot get image url: %s", err)
	}
	// get video
	videos , err := server.grpcClient.MediaClient.GetVideoUrl(ctx , &pb.GetUrlRequest{
		ProductId: product.ID.String(),
	})
	if err != nil {
		log.Error().Msgf("cannot get video url: %s", err)
	}

	return &pb.GetProductByIDResponse{
		Product: &pb.Product{
			Id:          product.ID.String(),
			CategoryId:  product.IDCategory.String(),
			IdAccount:   product.IDAccount.String(),
			Name:        product.Name,
			Description: product.Description,
			BrandName:   product.BrandName,
			Price:       product.Price,
			Inventory:   product.Inventory,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
		Images: &pb.Image{
			Url: image.Url,
		},
		Videos: &pb.Video{
			Url: videos.Url,
		},
		
	}, nil

}
