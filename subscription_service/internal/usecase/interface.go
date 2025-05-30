package usecase

import (
	"context"

	"github.com/recktt77/JobFree/subscription_service/internal/model"
)

type SubscriptionUsecase interface {
	// Подписки
	Subscribe(ctx context.Context, userID, planID string) (*model.Subscription, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	GetUserSubscriptions(ctx context.Context, userID string) ([]model.Subscription, error)
	GetSubscriptionStatus(ctx context.Context, userID string) (*model.Subscription, error)

	// Планы
	CreatePlan(ctx context.Context, plan model.Plan) error
	UpdatePlan(ctx context.Context, id string, update model.PlanUpdate) error
	GetPlanByID(ctx context.Context, id string) (*model.Plan, error)
	GetListOfPlans(ctx context.Context) ([]model.Plan, error)
}
