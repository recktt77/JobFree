package inmemory

import (
	"sync"

	"github.com/recktt77/JobFree/subscription_service/internal/model"
)

type Subscription struct {
	mu      sync.RWMutex
	storage map[string]model.Subscription
}

func NewSubscription() *Subscription {
	return &Subscription{
		storage: make(map[string]model.Subscription),
	}
}

func (s *Subscription) Set(sub model.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[sub.ID.Hex()] = sub
}

func (s *Subscription) SetMany(subs []model.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sub := range subs {
		s.storage[sub.ID.Hex()] = sub
	}
}

func (s *Subscription) Get(id string) (model.Subscription, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.storage[id]
	return val, ok
}

func (s *Subscription) GetAll() []model.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Subscription, 0, len(s.storage))
	for _, sub := range s.storage {
		result = append(result, sub)
	}
	return result
}
