package mongo

import (
	"context"
	"github.com/recktt77/JobFree/subscription_service/internal/model"
)

type Repository interface {
	// Plan
	CreatePlan(ctx context.Context, plan model.Plan) error
	GetPlanByID(ctx context.Context, filter model.PlanFilter) (model.Plan, error)
	UpdatePlan(ctx context.Context, filter model.PlanFilter, update model.PlanUpdate) error
	GetListOfPlans(ctx context.Context) ([]model.Plan, error)

	// Subscription
	CreateSubscription(ctx context.Context, sub model.Subscription) error
	CancelSubscription(ctx context.Context, id string) error
	GetUserSubscriptions(ctx context.Context, userID string) ([]model.Subscription, error)
	GetSubscriptionStatus(ctx context.Context, userID string) (*model.Subscription, error)
}
