package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetAccount(ctx context.Context, req *emptypb.Empty) (*pb.GetAccountResponse, error) {
	res, err := server.Auth(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get account: %s", err)
	}

	if !res.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get account: %s", "unauthorized")
	}

	account, err := server.store.GetAccount(ctx, uuid.MustParse(res.UserId))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "account not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get account: %s", err)
	}

	return &pb.GetAccountResponse{
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
