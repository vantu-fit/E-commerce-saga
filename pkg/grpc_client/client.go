package grpcclient

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	AccountClient pb.ServiceAccountClient
	ProductClient pb.ServiceProductClient
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) RunAccountClient(address string, doneCh chan struct{}) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		log.Info().Msg("Account client is closed")
	}()
	log.Info().Msg("Account client is running to connect to : " + address)
	c.AccountClient = pb.NewServiceAccountClient(conn)
	ping, err := c.AccountClient.Ping(context.Background(), &pb.PingRequest{
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

func (c *Client) RunProductClient(address string, doneCh chan struct{}) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		log.Info().Msg("Product client is closed")
	}()
	log.Info().Msg("Product client is running to connect to : " + address)
	c.ProductClient = pb.NewServiceProductClient(conn)
	ping, err := c.ProductClient.Ping(context.Background(), &pb.PingRequest{
		Message: "ping",
	})
	if err != nil {
		log.Error().Msgf("Product client can not ping: %v", err)
		return err
	}
	log.Info().Msg("Product client ping: " + ping.Message)
	<-doneCh
	return nil
}
