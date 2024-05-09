package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        interface{} `mapstructure:"service"`
	HTTP       HTTP        `mapstructure:"http"`
	Postgres   Postgres    `mapstructure:"postgres"`
	GRPC       GRPC        `mapstructure:"grpc"`
	GRPCClient GRPCClient  `mapstructure:"grpcClient"`
	Minio      Minio       `mapstructure:"minio"`
	Kafka      Kafka       `mapstructure:"kafka"`
}

type HTTP struct {
	Port string
	Mode string
}

type GRPC struct {
	Port              string
	Time              time.Duration
	Timeout           time.Duration
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	ShutdownWait      time.Duration
}

type GRPCClient struct {
	Timeout time.Duration
	Account string
	Product string
	Order   string
}

type Postgres struct {
	DnsURL string `mapstructure:"DNS_URL"`
	Migration string
}

type Minio struct {
	Endpoint string
	Username string
	Password string
}

type Kafka struct {
	Brokers []string
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
