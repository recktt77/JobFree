package events

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func NewNatsConn(url string) *nats.Conn {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	nc, err := nats.Connect(natsURL)

	if err != nil {
		log.Fatalf("NATS connection error: %v", err)
	}
	return nc
}
