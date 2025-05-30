package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host         string
	Password     string
	TLSEnable    bool
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type Subscription struct {
	subscription *redis.Client 
}

func NewProduct(ctx context.Context, cfg Config) (*Subscription, error) {

	opts := &redis.Options{
		Addr:         cfg.Host,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Password:     cfg.Password,
	}

	if cfg.TLSEnable {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	subscription := redis.NewClient(opts)

	err := subscription.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Subscription{subscription: subscription}, nil
}

func (p *Subscription) Close() error {
	err := p.subscription.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func (p *Subscription) Ping(ctx context.Context) error {
	err := p.subscription.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (p *Subscription) Unwrap() *redis.Client {
	return p.subscription
}