package inmemory

import (
	"context"
	_"errors"
	"sync"

	"github.com/recktt77/JobFree/payment_service/internal/model"
)

type PaymentRepo struct {
	data map[string]model.Payment
	mu   sync.RWMutex
}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{
		data: make(map[string]model.Payment),
	}
}

func (r *PaymentRepo) CreatePayment(ctx context.Context, payment model.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := payment.ID.Hex()
	r.data[id] = payment
	return nil
}

func (r *PaymentRepo) GetPaymentByID(ctx context.Context, id string) (*model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (r *PaymentRepo) ListPaymentsByUser(ctx context.Context, userID string) ([]model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Payment
	for _, p := range r.data {
		if p.UserID.Hex() == userID {
			result = append(result, p)
		}
	}
	return result, nil
}
