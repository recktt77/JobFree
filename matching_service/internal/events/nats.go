package events

import (
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNATS() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	log.Println("✅ Connected to NATS")
	return nc
}
