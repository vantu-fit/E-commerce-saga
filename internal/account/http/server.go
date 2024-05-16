package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/account/grpc"
	"github.com/vantu-fit/saga-pattern/pb"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

type HTTPGatewayServer struct {
	config      *config.Config
	store       db.Store
	grpcServer  *grpc.Server
	httpGateway *http.Server
}

func NewHTTPGatewayServer(
	config *config.Config,
	store db.Store,
	producer kafkaClient.Producer,
) (*HTTPGatewayServer, error) {
	auth := NewOAuth(config.Oauth.ClientID, config.Oauth.ClientSecret)
	ctx := context.Background()
	var err error
	server := &HTTPGatewayServer{
		config: config,
		store:  store,
	}

	server.grpcServer, err = grpc.NewServer(config, store, producer)
	if err != nil {
		return nil, err
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	err = pb.RegisterServiceAccountHandlerServer(ctx, grpcMux, server.grpcServer)
	if err != nil {
		return nil, errors.New("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	mux.HandleFunc("/api/v1/account/google/login", auth.Login)

	server.httpGateway = &http.Server{
		Handler: logger.HttpLogger(mux),
		Addr:    ":" + server.config.HTTP.Port,
	}

	return server, nil
}

func (server *HTTPGatewayServer) Run() error {
	return server.httpGateway.ListenAndServe()
}

func (server *HTTPGatewayServer) Shutdown(ctx context.Context) error {
	return server.httpGateway.Shutdown(ctx)
}
