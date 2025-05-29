package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"review_service/internal/model"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	conn *nats.Conn
}

func NewNatsBroker() (*NatsBroker, error) {
	url := os.Getenv("NATS_URL")
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}
	return &NatsBroker{conn: nc}, nil
}

func (b *NatsBroker) PublishReviewCreated(review *model.Review) {
	b.publish("review.created", review)
}

func (b *NatsBroker) PublishReviewModerated(reviewID string, action string) {
	msg := map[string]string{"review_id": reviewID, "action": action}
	b.publish("review.moderated", msg)
}

func (b *NatsBroker) publish(subject string, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal message for %s: %v", subject, err)
		return
	}
	if err := b.conn.Publish(subject, payload); err != nil {
		log.Printf("failed to publish %s: %v", subject, err)
	}
	log.Printf("[NATS PUBLISH] subject: %s → payload: %s\n", subject, string(payload))
}
