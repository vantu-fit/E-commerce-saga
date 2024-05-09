package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/orchestrator/config"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/event"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service"
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
	log.Info().Msg("Start orchestrator service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/orchestrator/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	// create context for gracefull shutdown
	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	// create kafka
	producer := kafkaClient.NewProducer(cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers)

	// create orchestrator service
	orchestratorService := service.NewService(producer)

	// create event handler
	orchestratorEventHandler := event.NewEvantHandler(cfg, consumer, orchestratorService)

	orchestratorEventHandler.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(0 * time.Second)
		log.Fatal().Msg("Graceful shutdown timeout")

		doneCh <- struct{}{}
	}()

	log.Info().Msg("product service shutdown")
	<-doneCh

}
