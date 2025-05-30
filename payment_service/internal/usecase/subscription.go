package usecase

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/payment_service/internal/adapter/nats/producer"
	"github.com/recktt77/JobFree/payment_service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentUsecase struct {
	repo     PaymentRepository
	producer *producer.PaymentProducer
}

func NewPaymentUsecase(repo PaymentRepository, producer *producer.PaymentProducer) *PaymentUsecase {
	return &PaymentUsecase{
		repo:     repo,
		producer: producer,
	}
}


func (u *PaymentUsecase) CreatePayment(ctx context.Context, payment model.Payment) error {
	payment.ID = primitive.NewObjectID()
	payment.Status = "pending"
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	if err := u.repo.CreatePayment(ctx, payment); err != nil {
		return err
	}

	_ = u.producer.PublishPaymentCreated(ctx, payment)

	return nil
}


func (u *PaymentUsecase) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	return u.repo.GetPayment(ctx, id)
}

func (u *PaymentUsecase) ListUserPayments(ctx context.Context, userID string) ([]model.Payment, error) {
	return u.repo.ListUserPayments(ctx, userID)
}
