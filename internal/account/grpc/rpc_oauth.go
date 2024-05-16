package grpc

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/hash"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Oauth(ctx context.Context, req *pb.OauthRequest) (*pb.LoginResponse, error) {
	auth := oauth2.Config{
		RedirectURL:  "http://localhost/api/v1/account/google/callback",
		ClientID:     s.config.Oauth.ClientID, 
		ClientSecret: s.config.Oauth.ClientSecret,                                      
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	code := req.GetCode()

	token, err := auth.Exchange(context.Background(), code)

	if err != nil {
		log.Error().Err(err).Msg("Failed to exchange token")
		return nil, err
	}

	client := auth.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user info")
		return nil, err
	}

	defer resp.Body.Close()

	var user pb.OauthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Error().Err(err).Msg("Failed to decode user info")
		return nil, err
	}

	// get account
	account, err := s.store.GetAccountByEmail(ctx, user.Email)
	if err == pgx.ErrNoRows {
		// create account
		hassPassword, err := hash.HashedPassword(user.Email)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
		}
		account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
			FirstName:   user.GivenName,
			LastName:    user.FamilyName,
			Email:       user.Email,
			Address:     "",
			PhoneNumber: randomPhoneNumber(),
			Password:    hassPassword,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create account")
			return nil, err
		}

		// create refreshtoken
		refreshToken, refreshPayload, err := s.maker.CreateToken(uuid.New(), account.ID, time.Hour*time.Duration(s.config.PasetoConfig.RefreshTokenExpire))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
		}

		// create accesstoken
		accessToken, _, err := s.maker.CreateToken(refreshPayload.ID, account.ID, time.Minute*time.Duration(s.config.PasetoConfig.AccessTokenExpire))
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

		session, err := s.store.CreateSession(ctx, argSession)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot create session: %s", err)
		}

		return &pb.LoginResponse{
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
		}, nil

	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to get account")
		return nil, err
	}

	// create refreshtoken
	refreshToken, refreshPayload, err := s.maker.CreateToken(uuid.New(), account.ID, time.Hour*time.Duration(s.config.PasetoConfig.RefreshTokenExpire))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
	}

	// create accesstoken
	accessToken, _, err := s.maker.CreateToken(refreshPayload.ID, account.ID, time.Minute*time.Duration(s.config.PasetoConfig.AccessTokenExpire))
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

	session, err := s.store.CreateSession(ctx, argSession)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %s", err)
	}

	return &pb.LoginResponse{
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
	}, nil
}

func randomPhoneNumber() string {
	n := 10
	var letter = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
