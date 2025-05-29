package producer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/JobFree/internal/adapter/nats/producer/dto"
	"github.com/recktt77/JobFree/internal/model"
)

type UserEventPublisher struct {
	nc *nats.Conn
}

func NewUserEventPublisher(nc *nats.Conn) *UserEventPublisher {
	return &UserEventPublisher{nc: nc}
}

func (p *UserEventPublisher) PublishUserRegistered(user *model.User) error {
	log.Println(">>> PublishUserRegistered called")
	event := dto.UserRegistered{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Profile.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.nc.Publish("user.registered", data)
	if err == nil {
		log.Printf("Published user.registered for %s", user.Email)
	}
	return err
}

func (p *UserEventPublisher) PublishProfileUpdated(userID string, profile model.UserProfile) error {
	event := dto.ProfileUpdated{
		ID:        userID,
		Profile:   profile,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.nc.Publish("user.profile_updated", data)
}
