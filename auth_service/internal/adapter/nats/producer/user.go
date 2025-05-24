package producer

import (
    "github.com/recktt77/JobFree/internal/model"
    "github.com/recktt77/JobFree/internal/adapter/nats/producer/dto"
    "encoding/json"
    "github.com/nats-io/nats.go"
    "time"
)

type UserEventPublisher struct {
    nc *nats.Conn
}

func NewUserEventPublisher(nc *nats.Conn) *UserEventPublisher {
    return &UserEventPublisher{nc: nc}
}

func (p *UserEventPublisher) PublishUserRegistered(user *model.User) error {
    event := dto.UserRegistered{
        ID:        user.ID.Hex(),
        Email:     user.Email,
        Role:      user.Role,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.nc.Publish("user.registered", data)
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
