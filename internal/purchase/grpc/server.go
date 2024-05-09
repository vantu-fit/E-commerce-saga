package grpc

import (
	"net"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/vantu-fit/saga-pattern/cmd/purchase/config"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service"
	"github.com/vantu-fit/saga-pattern/pb"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedServicePurchaseServer
	config     *config.Config
	service    *service.Service
	grpcServer *grpc.Server
	grpcClient *grpcclient.Client
}

func NewServer(
	config *config.Config,
	service *service.Service,
	grpcClient *grpcclient.Client,
) *Server {
	server := &Server{
		config:     config,
		service:    service,
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

	pb.RegisterServicePurchaseServer(server.grpcServer, server)

	reflection.Register(server.grpcServer)

	return server
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
