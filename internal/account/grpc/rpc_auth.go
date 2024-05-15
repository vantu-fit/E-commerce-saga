package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) Auth(ctx context.Context, req *emptypb.Empty) (*pb.AuthResponse, error) {
	var resoponse = &pb.AuthResponse{
		Valid:  false,
		UserId: "",
	}
	// check access token
	payload, err := server.authorizationUser(ctx)
	if err != nil {
		return resoponse, unauthenticatedError(err)
	}

	// check session
	session, err := server.store.GetSessionById(ctx, payload.ID)
	if err != nil {
		return resoponse, unauthenticatedError(err)
	}

	resoponse.Valid = true
	resoponse.UserId = session.UserID.String()

	return resoponse, nil
}
