package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/comment/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
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

	// user id
	userID := authResponse.UserId

	// create comment
	err = s.service.Command.CreateComment.Handle(ctx, command.CreateComment{
		CreateCommentRequest: req,
		UserID:               uuid.MustParse(userID),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal , "Internal error")
	}

	return &pb.CreateCommentResponse{
		Comment: &pb.Comment{
			Content: req.GetContent(),
		},
	}, nil

}
