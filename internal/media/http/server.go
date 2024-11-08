package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/grpc"
	"github.com/vantu-fit/saga-pattern/internal/media/service"
	"github.com/vantu-fit/saga-pattern/pb"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

type HTTPGatewayServer struct {
	config      *config.Config
	store       db.Store
	service     *service.Service
	grpcServer  *grpc.Server
	httpGateway *http.Server
}

func NewHTTPGatewayServer(
	config *config.Config,
	store db.Store,
	service *service.Service,
	grpcClient *grpcclient.Client,
) (*HTTPGatewayServer, error) {
	ctx := context.Background()
	var err error
	server := &HTTPGatewayServer{
		config:  config,
		service: service,
		store:   store,
	}

	server.grpcServer = grpc.NewServer(config, store, service, grpcClient)

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

	err = pb.RegisterServiceMediaHandlerServer(ctx, grpcMux, server.grpcServer)
	if err != nil {
		return nil, errors.New("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	mux.Handle("/api/v1/media/image/upload", &imageUploadHandler{handler: server.ImageUploadHandler})

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
