package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) GetClient() *redis.Client {
	return r.client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: rdb}
}

// Сохраняем фрилансера в Redis
func (r *RedisCache) SaveFreelancer(freelancerID string, data map[string]interface{}) error {
	key := fmt.Sprintf("freelancer:%s", freelancerID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonData, 0).Err()
}

// Получаем всех фрилансеров
func (r *RedisCache) GetAllFreelancers() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	iter := r.client.Scan(ctx, 0, "freelancer:*", 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.client.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var f map[string]interface{}
		if err := json.Unmarshal([]byte(val), &f); err != nil {
			continue
		}
		result = append(result, f)
	}

	return result, nil
}
