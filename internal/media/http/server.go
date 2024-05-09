package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/grpc"
	"github.com/vantu-fit/saga-pattern/internal/media/service"
	"github.com/vantu-fit/saga-pattern/internal/media/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

type HTTPGatewayServer struct {
	config      *config.Config
	store       db.Store
	grpcServer  *grpc.Server
	httpGateway *http.Server
}

func NewHTTPGatewayServer(config *config.Config, store db.Store, service *service.Service, grpcClient *grpcclient.Client) (*HTTPGatewayServer, error) {
	ctx := context.Background()
	var err error
	server := &HTTPGatewayServer{
		config: config,
		store:  store,
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
	mux.Handle("/v1/media/image", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Auth
		// validate request
		file, _, err := r.FormFile("data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Call command
		url, err := service.Commad.UploadImage.Handle(ctx, command.UploadImage{
			UploadImageRequest: &pb.UploadImageRequest{
				Filename:  r.FormValue("filename"),
				Alt:       r.FormValue("alt"),
				Data:      buf.Bytes(),
				ProductId: r.FormValue("product_id"),
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"url":"` + url + `"}`))
	}))

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
