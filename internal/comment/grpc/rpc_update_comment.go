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

func (s *Server) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
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

	// update comment
	err = s.service.Command.UpdateComment.Handle(ctx, command.UpdateComment{
		UpdateCommentRequest: req,
		UserID:               uuid.MustParse(userID),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error , err: %v", err)
	}

	return &pb.UpdateCommentResponse{
		Comment: &pb.Comment{
			Id:      req.GetId(),
			Content: req.GetContent(),
		},
	}, nil
}
