package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/recktt77/JobFree/subscription_service/internal/model"
	"github.com/recktt77/JobFree/subscription_service/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const keyPrefix = "subscription:user:%s"

type Subscription struct {
	rdb *redis.Subscription
	ttl time.Duration
}

func NewSubscription(rdb *redis.Subscription, ttl time.Duration) *Subscription {
	return &Subscription{
		rdb: rdb,
		ttl: ttl,
	}
}

func (s *Subscription) Set(ctx context.Context, sub model.Subscription) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}
	fmt.Println("Set in Redis:", sub.ID.Hex())
	return s.rdb.Unwrap().Set(ctx, s.key(sub.ID), data, s.ttl).Err()
}

func (s *Subscription) SetMany(ctx context.Context, subs []model.Subscription) error {
	pipe := s.rdb.Unwrap().Pipeline()
	for _, sub := range subs {
		data, err := json.Marshal(sub)
		if err != nil {
			return fmt.Errorf("failed to marshal subscription: %w", err)
		}
		pipe.Set(ctx, s.key(sub.ID), data, s.ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set many subscriptions: %w", err)
	}
	return nil
}

func (s *Subscription) Get(ctx context.Context, id string) (model.Subscription, error) {
	fmt.Println("Try get from Redis:", id)
	data, err := s.rdb.Unwrap().Get(ctx, fmt.Sprintf(keyPrefix, id)).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return model.Subscription{}, nil
		}
		return model.Subscription{}, fmt.Errorf("failed to get subscription: %w", err)
	}

	var sub model.Subscription
	err = json.Unmarshal(data, &sub)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return sub, nil
}

func (s *Subscription) GetAll(ctx context.Context) ([]model.Subscription, error) {
	keys, err := s.rdb.Unwrap().Keys(ctx, "subscription:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscriptions: %w", err)
	}

	var subs []model.Subscription
	for _, key := range keys {
		data, err := s.rdb.Unwrap().Get(ctx, key).Bytes()
		if err != nil {
			if err == goredis.Nil {
				continue
			}
			return nil, fmt.Errorf("failed to get subscription: %w", err)
		}

		var sub model.Subscription
		err = json.Unmarshal(data, &sub)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (s *Subscription) key(id primitive.ObjectID) string {
	return fmt.Sprintf(keyPrefix, id.Hex())
}
