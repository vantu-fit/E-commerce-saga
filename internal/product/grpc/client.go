package grpc

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config        *config.Config
	accountClient pb.ServiceAccountClient
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) RunAccountClient(doneCh chan struct{}) error {
	conn, err := grpc.Dial(c.config.GRPCClient.Account, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func () {
		conn.Close()
		log.Info().Msg("Account client is closed")
	}()
	log.Info().Msg("Account client is running to connect to : " + c.config.GRPCClient.Account)
	c.accountClient = pb.NewServiceAccountClient(conn)
	ping, err := c.accountClient.Ping(context.Background(), &pb.PingRequest{
		Message: "ping",
	})
	if err != nil {
		log.Error().Msgf("Account client can not ping: %v", err)
		return err
	}
	log.Info().Msg("Account client ping: " + ping.Message)
	<-doneCh
	return nil
}
