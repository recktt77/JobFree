package producer


import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/recktt77/JobFree/payment_service/internal/adapter/nats/producer/dto"
	"github.com/recktt77/JobFree/payment_service/internal/model"
)

type PaymentProducer struct {
	conn *nats.Conn
}

func NewPaymentProducer(conn *nats.Conn) *PaymentProducer {
	return &PaymentProducer{conn: conn}
}

func (p *PaymentProducer) PublishPaymentCreated(ctx context.Context, payment model.Payment) error {
	event := dto.PaymentCreated{
		ID:             payment.ID.Hex(),
		UserID:         payment.UserID.Hex(),
		SubscriptionID: payment.SubscriptionID.Hex(),
		Amount:         payment.Amount,
		Status:         payment.Status,
		CreatedAt:      payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.conn.Publish("payment.created", data)
}
