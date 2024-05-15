package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vantu-fit/saga-pattern/cmd/mail/config"
	"github.com/vantu-fit/saga-pattern/internal/mail/event"
	"github.com/vantu-fit/saga-pattern/internal/mail/sender"
	"github.com/vantu-fit/saga-pattern/internal/mail/service"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// load config
	cfgFile, err := config.LoadConfig("./cmd/mail/config/config")
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

	// create kafka producer
	producer := kafkaClient.NewProducer(cfg.Kafka.Brokers)
	consummer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers)

	// create mail builder
	builder := sender.NewMailBuilder().
		SetDomain(cfg.Mail.MailDomain).
		SetHost(cfg.Mail.MailHostSend).
		SetPort(cfg.Mail.MailPortSend).
		SetUsername(cfg.Mail.MailUsername).
		SetPassword(cfg.Mail.MailPassword).
		SetEncryption(cfg.Mail.MailEncryption).
		SetFromAddress(cfg.Mail.MailFromAddress).
		SetFromName(cfg.Mail.MailFromName)

	mail := builder.Build()

	// create mail service
	mailService := service.NewService(mail)

	// create event handler
	eventHanlder := event.NewEvantHandler(cfg, consummer, producer, mailService)

	// run event handler
	eventHanlder.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(0 * time.Second)

		doneCh <- struct{}{}
	}()

	<-doneCh
	log.Info().Msg("Account service shutdown")
}
