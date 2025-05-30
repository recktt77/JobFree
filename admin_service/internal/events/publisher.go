package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

type EventPublisher struct {
	nc *nats.Conn
}

func NewEventPublisher(nc *nats.Conn) *EventPublisher {
	return &EventPublisher{nc: nc}
}

func (p *EventPublisher) PublishUserBanned(ctx context.Context, userID, reason string) error {
	event := UserBannedEvent{
		UserID:    userID,
		Reason:    reason,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish("ap2.admin.user.banned", data)
}
