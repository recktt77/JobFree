package mongo

import (
	"context"

	"github.com/recktt77/JobFree/payment_service/internal/adapter/mongo/dao"
	"github.com/recktt77/JobFree/payment_service/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type payment struct {
	dao PaymentRepository
}

func NewPaymentRepository(db *mongo.Database) PaymentRepository {
	return &payment{
		dao: dao.NewPaymentRepository(db),
	}
}

func (p *payment) CreatePayment(ctx context.Context, payment model.Payment) error {
	return p.dao.CreatePayment(ctx, payment)
}

func (p *payment) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	return p.dao.GetPayment(ctx, id)
}

func (p *payment) ListUserPayments(ctx context.Context, userID string) ([]model.Payment, error) {
	return p.dao.ListUserPayments(ctx, userID)
}
