package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App   any `mapstructure:"service"`
	HTTP  HTTP `mapstructure:"http"`
	GRPC  GRPC `mapstructure:"grpc"`
	Kafka Kafka `mapstructure:"kafka"`
	Mail  Mail `mapstructure:"mail"`
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
	Product string
	Account string
}

type Kafka struct {
	Brokers []string
}

type Mail struct {
	MailDomain      string
	MailHostSend    string
	MailPortSend    int
	MailUsername    string
	MailPassword    string
	MailEncryption  string
	MailFromName    string
	MailFromAddress string
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
