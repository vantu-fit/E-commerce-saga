package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/vantu-fit/saga-pattern/internal/account/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader     = "Authorization"
	authorizationTypeBearer = "Bearer"
)

func (server *Server) authorizationUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization format")
	}

	authType := fields[0]
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupport authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.maker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthenticated error: %s", err)
}
