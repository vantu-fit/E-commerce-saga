package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/hash"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := validateLoginRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// get account by email
	account, err := server.store.GetAccountByEmail(ctx, req.GetEmail())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "account not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get account by email: %s", err)
	}

	// compare password
	err = hash.CheckPassword(account.Password, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %s", err)
	}
	// create refreshtoken
	refreshToken, refreshPayload, err := server.maker.CreateToken(uuid.New(), account.ID, time.Hour*time.Duration(server.config.PasetoConfig.RefreshTokenExpire))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
	}

	// create accesstoken
	accessToken, _, err := server.maker.CreateToken(refreshPayload.ID, account.ID, time.Minute*time.Duration(server.config.PasetoConfig.AccessTokenExpire))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	// create session
	argSession := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       refreshPayload.UserID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
	}

	session, err := server.store.CreateSession(ctx, argSession)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %s", err)
	}

	response := pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionId:    session.ID.String(),
		Account: &pb.Account{
			Id:          account.ID.String(),
			FirstName:   account.FirstName,
			LastName:    account.LastName,
			Email:       account.Email,
			PhoneNumber: account.PhoneNumber,
			Address:     account.Address,
			Active:      account.Active.Bool,
			UpdatedAt:   timestamppb.New(account.UpdatedAt),
			CreatedAt:   timestamppb.New(account.CreatedAt),
		},
	}

	return &response, nil
}

func validateLoginRequest(req *pb.LoginRequest) (violatios []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violatios = append(violatios, &errdetails.BadRequest_FieldViolation{
			Field:       "email",
			Description: err.Error(),
		})
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violatios = append(violatios, &errdetails.BadRequest_FieldViolation{
			Field:       "password",
			Description: err.Error(),
		})
	}

	return violatios
}
