package usecase

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/subscription_service/internal/adapter/mongo"
	"github.com/recktt77/JobFree/subscription_service/internal/adapter/nats/producer"
	dto "github.com/recktt77/JobFree/subscription_service/internal/adapter/nats/producer/dto"
	"github.com/recktt77/JobFree/subscription_service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type subscriptionUsecase struct {
	repo     mongo.Repository
	producer *producer.Producer
}

func NewSubscriptionUsecase(repo mongo.Repository, producer *producer.Producer) *subscriptionUsecase {
	return &subscriptionUsecase{
		repo:     repo,
		producer: producer,
	}
}

// === PLAN ===

func (u *subscriptionUsecase) CreatePlan(ctx context.Context, plan model.Plan) error {
	return u.repo.CreatePlan(ctx, plan)
}

func (u *subscriptionUsecase) UpdatePlan(ctx context.Context, id string, update model.PlanUpdate) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := model.PlanFilter{ID: &objID}
	return u.repo.UpdatePlan(ctx, filter, update)
}

func (u *subscriptionUsecase) GetPlanByID(ctx context.Context, id string) (*model.Plan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := model.PlanFilter{ID: &objID}
	plan, err := u.repo.GetPlanByID(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (u *subscriptionUsecase) GetListOfPlans(ctx context.Context) ([]model.Plan, error) {
	return u.repo.GetListOfPlans(ctx)
}

// === SUBSCRIPTION ===

func (u *subscriptionUsecase) Subscribe(ctx context.Context, userID, planID string) (*model.Subscription, error) {
	start := time.Now()
	end := start.AddDate(0, 0, 30)

	sub := model.Subscription{
		ID:        model.NewObjectID(),
		UserID:    model.MustObjectIDFromHex(userID),
		PlanID:    model.MustObjectIDFromHex(planID),
		StartDate: start,
		EndDate:   end,
		Active:    true,
	}

	if err := u.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	if u.producer != nil {
		payload := dto.ToSubscriptionPayload(sub)
		_ = u.producer.PublishEvent("subscription_service.events", "subscription.created", payload)
	}

	return &sub, nil
}

func (u *subscriptionUsecase) CancelSubscription(ctx context.Context, subscriptionID string) error {
	err := u.repo.CancelSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if u.producer != nil {
		sub := model.Subscription{
			ID:     model.MustObjectIDFromHex(subscriptionID),
			Active: false,
		}
		payload := dto.ToSubscriptionPayload(sub)
		_ = u.producer.PublishEvent("subscription_service.events", "subscription.cancelled", payload)
	}

	return nil
}

func (u *subscriptionUsecase) GetUserSubscriptions(ctx context.Context, userID string) ([]model.Subscription, error) {
	return u.repo.GetUserSubscriptions(ctx, userID)
}

func (u *subscriptionUsecase) GetSubscriptionStatus(ctx context.Context, userID string) (*model.Subscription, error) {
	return u.repo.GetSubscriptionStatus(ctx, userID)
}
