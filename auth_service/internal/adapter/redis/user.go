package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/recktt77/JobFree/internal/model"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func (r *RedisCache) Close() any {
	panic("unimplemented")
}

func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *RedisCache) GetUser(ctx context.Context, id string) (*model.User, error) {
	key := fmt.Sprintf("user:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisCache) SetUser(ctx context.Context, user *model.User) error {
	key := fmt.Sprintf("user:%s", user.ID.Hex())
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}
