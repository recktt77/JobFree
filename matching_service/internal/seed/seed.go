package seed

import (
	"log"

	"github.com/recktt77/JobFree/matching_service/internal/cache"
)

func Run(redis *cache.RedisCache) {
	freelancers := []map[string]interface{}{
		{
			"freelancer_id": "f123",
			"name":          "Test User",
			"skills":        []string{"Go", "MongoDB"},
		},
		{
			"freelancer_id": "f456",
			"name":          "Backend Pro",
			"skills":        []string{"Go", "Redis", "Docker"},
		},
	}

	for _, f := range freelancers {
		id := f["freelancer_id"].(string)
		err := redis.SaveFreelancer(id, f)
		if err != nil {
			log.Printf("failed to insert %s to Redis: %v", id, err)
		} else {
			log.Printf("seeded freelancer %s", id)
		}
	}
}
