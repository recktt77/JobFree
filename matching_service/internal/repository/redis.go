package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(addr string) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // "localhost:6379"
		DB:   0,
	})

	return &RedisRepository{Client: rdb}
}

func (r *RedisRepository) SaveDummyBid(ctx context.Context, bidID string) error {
	return r.Client.Set(ctx, fmt.Sprintf("bid:%s", bidID), "dummy", 0).Err()
}
