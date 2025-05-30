package usecase

import (
	"context"
	"github.com/recktt77/JobFree/payment_service/internal/model"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment model.Payment) error
	GetPayment(ctx context.Context, id string) (*model.Payment, error)
	ListUserPayments(ctx context.Context, userID string) ([]model.Payment, error)
}
