package grpc

import (
	"log"
	"net"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/account/token"
	"github.com/vantu-fit/saga-pattern/pb"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
	"github.com/vantu-fit/saga-pattern/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedServiceAccountServer
	config     *config.Config
	maker      token.Maker
	store      db.Store
	producer   kafkaClient.Producer
	grpcServer *grpc.Server
}

func NewServer(
	config *config.Config, 
	store db.Store ,
	producer kafkaClient.Producer,
) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.PasetoConfig.SymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config: config,
		store:  store,
		maker:  maker,
		producer: producer,
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

	pb.RegisterServiceAccountServer(server.grpcServer, server)

	reflection.Register(server.grpcServer)

	return server, nil
}

func (server *Server) Run() error {
	addr := "0.0.0.0:" + server.config.GRPC.Port
	log.Println("grpc server listening on ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if err := server.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (srv *Server) GracefulStop() {
	srv.grpcServer.GracefulStop()
}
