package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"messaging_service/email"
	"messaging_service/internal/repository"

	"github.com/nats-io/nats.go"
)

type UserRegisteredEvent struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type NatsCoreSubscriber struct {
	nc       *nats.Conn
	repo     repository.MessageRepository
	notifier email.EmailSender
}

func NewNatsSubscriberCore(nc *nats.Conn, repo repository.MessageRepository, notifier email.EmailSender) *NatsCoreSubscriber {
	return &NatsCoreSubscriber{nc, repo, notifier}
}

func (s *NatsCoreSubscriber) SubscribeToUserRegistered() error {
	_, err := s.nc.Subscribe("user.registered", func(msg *nats.Msg) {
		var evt UserRegisteredEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Println("failed to unmarshal event:", err)
			return
		}

		log.Printf("Received user.registered: %v", evt.Email)

		// Send email
		notificationText := fmt.Sprintf("Welcome to JobFree, %s!", evt.Name)
		err := s.notifier.Send(evt.Email, "Welcome!", notificationText)
		if err != nil {
			log.Printf("email send error: %v", err)
		}
	})

	return err
}
