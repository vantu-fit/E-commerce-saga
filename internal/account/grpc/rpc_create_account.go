package grpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
	"github.com/vantu-fit/saga-pattern/pkg/hash"
	val "github.com/vantu-fit/saga-pattern/pkg/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	violations := validateCreteUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hassPassword, err := hash.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateAccountParams{
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Address:     req.GetAddress(),
		Password:    hassPassword,
	}

	// create account
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create account: %s", err)
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
		UserID:       account.ID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
	}

	session, err := server.store.CreateSession(ctx, argSession)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %s", err)
	}

	response := pb.CreateAccountResponse{
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

	pdSendMail := pb.SendRegisterEmailRequest{
		FromName: "Saga Orchestrator",
		From:     "saga@gmail.com",
		To:       account.Email,
		Subject:  "Register Success",
	}

	payload, err := json.Marshal(&pdSendMail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot marshal payload: %s", err)
	}

	err = server.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.SendRegisterEmailTopic,
		Value: payload,
	})

	if err != nil {
		log.Error().Msgf("cannot publish message: %s", err)
	}

	return &response, nil
}

func validateCreteUserRequest(req *pb.CreateAccountRequest) (violatios []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetFirstName()); err != nil {
		violatios = append(violatios, fileViolation("first_name", err))
	}
	if err := val.ValidateFullname(req.GetLastName()); err != nil {
		violatios = append(violatios, fileViolation("last_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violatios = append(violatios, fileViolation("email", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violatios = append(violatios, fileViolation("password", err))
	}
	if err := val.ValidatePhoneNumber(req.GetPhoneNumber()); err != nil {
		violatios = append(violatios, fileViolation("phone_number", err))
	}
	if err := val.ValidateString(req.GetAddress(), 3, 100); err != nil {
		violatios = append(violatios, fileViolation("address", err))
	}

	return violatios
}

func fileViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusIvalid := status.New(codes.InvalidArgument, "invalid parameter")

	statusDetails, err := statusIvalid.WithDetails(badRequest)
	if err != nil {
		return statusIvalid.Err()
	}

	return statusDetails.Err()
}
