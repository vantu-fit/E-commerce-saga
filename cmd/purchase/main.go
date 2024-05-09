package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/purchase/config"
	"github.com/vantu-fit/saga-pattern/internal/purchase/event"
	"github.com/vantu-fit/saga-pattern/internal/purchase/grpc"
	"github.com/vantu-fit/saga-pattern/internal/purchase/http"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// log for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// start service
	log.Info().Msg("Start purchase service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/purchase/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}
	log.Info().Msgf("account : %s" , cfg.GRPCClient.Account)
	log.Info().Msgf("product : %s" , cfg.GRPCClient.Product)

	// create context for gracefull shutdown
	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	// create kafka client
	producer := kafkaClient.NewProducer(cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers)

	// create event handler
	eventHandler := event.NewEvantHandler(cfg, producer , consumer)

	// create purchase service
	purchaseService := service.NewService(eventHandler)

	// create grpc client
	grpcClient := grpcclient.NewClient()

	// run account grpc client
	go func() {
		if err := grpcClient.RunAccountClient(cfg.GRPCClient.Account, doneCh); err != nil {
			log.Fatal().Msgf("Run grpc account client: %v", err)
		}
	}()

	// run product grpc client
	go func() {
		if err := grpcClient.RunProductClient(cfg.GRPCClient.Product, doneCh); err != nil {
			log.Fatal().Msgf("Run grpc product client: %v", err)
		}
	}()

	// create grpc server
	grpcServer := grpc.NewServer(cfg, purchaseService, grpcClient)

	// run grpc server
	go func() {
		log.Info().Msgf("GRPC server is running on port %s", cfg.GRPC.Port)
		if err := grpcServer.Run(); err != nil {
			log.Fatal().Msgf("Run grpc server: %v", err)
		}
	}()

	// create http gateway server
	HTTPGatewayServer, err := http.NewHTTPGatewayServer(cfg, grpcClient, purchaseService)
	if err != nil {
		log.Fatal().Msgf("Create http gateway server: %v", err)
	}

	// run http gateway server
	go func() {
		log.Info().Msgf("HTTP Gateway server is running on port %s", cfg.HTTP.Port)
		if err := HTTPGatewayServer.Run(); err != nil {
			log.Fatal().Msgf("Run http gateway server: %v", err)
		}
	}()

	// run event handler
	eventHandler.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(0 * time.Second)
		log.Fatal().Msg("Graceful shutdown timeout")

		HTTPGatewayServer.Shutdown(context.Background())

		grpcServer.GracefulStop()
		doneCh <- struct{}{}
	}()

	log.Info().Msg("product service shutdown")
	<-doneCh

}
