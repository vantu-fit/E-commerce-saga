package grpc

import (
	"context"

	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Auth(ctx context.Context, req *emptypb.Empty) (*pb.AuthResponse, error) {
	payload, err := server.authorizationUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	account, err := server.store.GetAccountByEmail(ctx, payload.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get account: %s", err)
	}

	return &pb.AuthResponse{
		Valid: true,
		Account: &pb.Account{
			Id:          account.ID.String(),
			FirstName:   account.FirstName,
			LastName:    account.LastName,
			Email:       account.Email,
			Address:     account.Address,
			PhoneNumber: account.PhoneNumber,
			Active:      account.Active.Bool,
			UpdatedAt:   timestamppb.New(account.UpdatedAt),
			CreatedAt:   timestamppb.New(account.CreatedAt),
		},
	}, nil
}
