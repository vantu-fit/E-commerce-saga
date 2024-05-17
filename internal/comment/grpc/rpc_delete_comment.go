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

func (s *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
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
	err = s.service.Command.DeleteComment.Handle(ctx, command.DeleteComment{
		DeleteCommentRequest: req,
		UserID: 			 uuid.MustParse(userID),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal , "Internal error, err: %v", err)
	}

	return &pb.DeleteCommentResponse{
		Comment: &pb.Comment{
			Id: req.GetId(),
		},
	}, nil

}
