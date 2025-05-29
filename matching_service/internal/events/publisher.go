package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/JobFree/matching_service/internal/model"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(natsURL string) *Publisher {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}

	return &Publisher{conn: nc}
}

func (p *Publisher) PublishBidCreated(ctx context.Context, bid *model.Bid) error {
	data, err := json.Marshal(bid)
	if err != nil {
		return err
	}

	subject := "ap2.statistics.event.updated"
	return p.conn.Publish(subject, data)
}
