package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"review_service/internal/model"

	"github.com/redis/go-redis/v9"
)

type ReviewCache struct {
	client *redis.Client
}

func NewReviewCache() (*ReviewCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	return &ReviewCache{client: rdb}, nil
}

func keyByProject(projectID string) string {
	return fmt.Sprintf("review:project:%s", projectID)
}

func keyByReviewee(revieweeID string) string {
	return fmt.Sprintf("review:user:%s", revieweeID)
}

func (c *ReviewCache) SetProjectReviews(ctx context.Context, projectID string, reviews []model.Review) error {
	data, err := json.Marshal(reviews)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, keyByProject(projectID), data, 10*time.Minute).Err()
}

func (c *ReviewCache) SetUserReviews(ctx context.Context, revieweeID string, reviews []model.Review) error {
	data, err := json.Marshal(reviews)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, keyByReviewee(revieweeID), data, 10*time.Minute).Err()
}

func (c *ReviewCache) GetProjectReviews(ctx context.Context, projectID string) ([]model.Review, error) {
	var reviews []model.Review
	data, err := c.client.Get(ctx, keyByProject(projectID)).Bytes()
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(data, &reviews)
	return reviews, nil
}

func (c *ReviewCache) GetUserReviews(ctx context.Context, revieweeID string) ([]model.Review, error) {
	var reviews []model.Review
	data, err := c.client.Get(ctx, keyByReviewee(revieweeID)).Bytes()
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(data, &reviews)
	return reviews, nil
}

func (c *ReviewCache) Invalidate(ctx context.Context, projectID, revieweeID string) {
	if projectID != "" {
		key := keyByProject(projectID)
		err := c.client.Del(ctx, key).Err()
		if err != nil {
			fmt.Printf("[CACHE INVALIDATE ERROR] project_id: %s → %v\n", projectID, err)
		} else {
			fmt.Println("[CACHE INVALIDATE] project_id:", projectID)
		}
	}

	if revieweeID != "" {
		key := keyByReviewee(revieweeID)
		err := c.client.Del(ctx, key).Err()
		if err != nil {
			fmt.Printf("[CACHE INVALIDATE ERROR] reviewee_id: %s → %v\n", revieweeID, err)
		} else {
			fmt.Println("[CACHE INVALIDATE] reviewee_id:", revieweeID)
		}
	}
}
