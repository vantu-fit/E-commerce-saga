package grpc

import (
	"net"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedServiceProductServer
	config     *config.Config
	store      db.Store
	grpcServer *grpc.Server
	grpcClient *Client
}

func NewServer(config *config.Config, store db.Store, grpcClient *Client) (*Server, error) {
	server := &Server{
		config:     config,
		store:      store,
		grpcClient: grpcClient,
	}

	server.grpcServer = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: config.GRPC.MaxConnectionIdle * time.Second,
			MaxConnectionAge:  config.GRPC.MaxConnectionAge * time.Minute,
			Timeout:           config.GRPC.Timeout * time.Second,
			Time:              config.GRPC.Time * time.Second,
		}),
		grpc.ChainUnaryInterceptor(
			grpcrecovery.UnaryServerInterceptor(),
		),
		grpc.UnaryInterceptor(logger.GrpcLogger),
	)

	pb.RegisterServiceProductServer(server.grpcServer, server)

	reflection.Register(server.grpcServer)

	return server, nil
}

func (server *Server) Run() error {
	addr := "0.0.0.0:" + server.config.GRPC.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if err := server.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (server *Server) GracefulStop() {
	server.grpcServer.GracefulStop()
}
