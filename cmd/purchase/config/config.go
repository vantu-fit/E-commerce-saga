package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App         any
	HTTP        HTTP
	GRPC        GRPC
	GRPCClient  GRPCClient
	Postgres    Postgres
	Migration   Migration
	RpcEnpoints RpcEndpoints `mapstructure:"rpcEndpoints"`
	Kafka       Kafka
	LocalCache  LocalCache `mapstructure:"localCache"`
	RedisCache  RedisCache `mapstructure:"redisCache"`
}

type Service struct {
	Name         string
	Mode         string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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
}

type GRPCClient struct {
	Timeout time.Duration
	Account string
	Product string
	Order   string
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type Postgres struct {
	DnsURL string `mapstructure:"DNS_URL"`
}

type Migration struct {
	Enable   bool
	Recreate bool
}

type RpcEndpoints struct {
	AuthSvc string
}

type Kafka struct {
	Brokers []string
}

type LocalCache struct {
	ExpirationTime uint64
}

type RedisCache struct {
	Address        []string
	Password       string
	DB             int
	PoolSize       int
	MaxRetries     int
	ExpirationTime uint64
	CuckooFilter   CuckooFilter `mapstructure:"CuckooFilter"`
}

type CuckooFilter struct {
	Capacity      int64
	BucketSize    int64
	MaxIterations int64
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
