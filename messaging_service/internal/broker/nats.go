package broker

import (
	"encoding/json"
	"log"
	"messaging_service/email"
	"messaging_service/internal/model"
	"messaging_service/internal/repository"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegisteredEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
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

		userID, _ := primitive.ObjectIDFromHex(evt.UserID)
		notification := &model.Notification{
			UserID:    userID,
			Type:      "registration_success",
			Message:   "Welcome, " + evt.Name + "!",
			Read:      false,
			CreatedAt: time.Now(),
		}

		// сохранять некуда, если у тебя нет NotificationRepository — просто отправим email
		if err := s.notifier.Send(evt.Email, "Welcome to our platform!", notification.Message); err != nil {
			log.Println("email send failed:", err)
		} else {
			log.Println("welcome email sent to", evt.Email)
		}
	})
	return err
}
