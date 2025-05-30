package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/recktt77/JobFree/subscription_service/pkg/mongo"
)

type (
	Config struct {
		Mongo   mongo.Config
		Server  Server
		Version string `env:"VERSION"`

		NATSUrl string `env:"NATS_URL" envDefault:"nats://localhost:4222"`

		Redis RedisConfig
		Cache CacheConfig
	}

	Server struct {
		GRPCServer GRPCServer
	}

	GRPCServer struct {
		Port                  int           `env:"GRPC_PORT,notEmpty"`
		MaxRecvMsgSizeMiB     int           `env:"GRPC_MAX_MESSAGE_SIZE_MIB" envDefault:"12"`
		MaxConnectionAge      time.Duration `env:"GRPC_MAX_CONNECTION_AGE" envDefault:"30s"`
		MaxConnectionAgeGrace time.Duration `env:"GRPC_MAX_CONNECTION_AGE_GRACE" envDefault:"10s"`
	}

	RedisConfig struct {
		Host         string        `env:"REDIS_HOST,notEmpty"`
		Password     string        `env:"REDIS_PASSWORD"`
		TLSEnable    bool          `env:"REDIS_TLS_ENABLE" envDefault:"false"`
		DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"5s"`
		ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"5s"`
	}

	CacheConfig struct {
		SubscriptionTTL time.Duration `env:"SUBSCRIPTION_CACHE_TTL" envDefault:"24h"`
	}
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
