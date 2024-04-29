package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       interface{}
	HTTP      HTTP
	GRPC      GRPC
	Postgres  Postgres
	Migration Migration
	PasetoConfig PasetoConfig
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

type Postgres struct {
	DnsURL    string `mapstructure:"DNS_URL"`
	Migration string `mapstructure:"Migration"`
}

type Migration struct {
	Enable   bool
	Recreate bool
}

type PasetoConfig struct {
	SymmetricKey       string `mapstructure:"SymmetricKey"`
	AccessTokenExpire  uint
	RefreshTokenExpire uint
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