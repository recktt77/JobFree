package frontend

import (
	"context"

	"github.com/recktt77/JobFree/subscription_service/internal/adapter/grpc/server/frontend/dto"
	"github.com/recktt77/JobFree/subscription_service/internal/model"
	"github.com/recktt77/JobFree/subscription_service/internal/usecase"
	subpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Subscription struct {
	subpb.UnimplementedSubscriptionServiceServer
	uc usecase.SubscriptionUsecase
}

func NewSubscription(uc usecase.SubscriptionUsecase) *Subscription {
	return &Subscription{uc: uc}
}

func (s *Subscription) Subscribe(ctx context.Context, req *subpb.SubscribeRequest) (*subpb.SubscribeResponse, error) {
	sub, err := s.uc.Subscribe(ctx, req.GetUserId(), req.GetPlanId())
	if err != nil {
		return nil, err
	}
	return &subpb.SubscribeResponse{Subscription: dto.ToProto(*sub)}, nil
}

func (s *Subscription) CancelSubscription(ctx context.Context, req *subpb.CancelSubscriptionRequest) (*subpb.CancelSubscriptionResponse, error) {
	err := s.uc.CancelSubscription(ctx, req.GetSubscriptionId())
	if err != nil {
		return nil, err
	}
	return &subpb.CancelSubscriptionResponse{Message: "Subscription cancelled"}, nil
}

func (s *Subscription) GetSubscriptionStatus(ctx context.Context, req *subpb.GetSubscriptionStatusRequest) (*subpb.GetSubscriptionStatusResponse, error) {
	sub, err := s.uc.GetSubscriptionStatus(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &subpb.GetSubscriptionStatusResponse{Subscription: dto.ToProto(*sub)}, nil
}

func (s *Subscription) GetSubscription(ctx context.Context, req *subpb.GetSubscriptionRequest) (*subpb.GetSubscriptionResponse, error) {
	subs, err := s.uc.GetUserSubscriptions(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	var result []*subpb.Subscription
	for _, sub := range subs {
		result = append(result, dto.ToProto(sub))
	}
	return &subpb.GetSubscriptionResponse{Subscriptions: result}, nil
}



func (s *Subscription) CreatePlan(ctx context.Context, req *subpb.CreatePlanRequest) (*subpb.CreatePlanResponse, error) {
	plan := dto.FromProtoToPlan(req)
	err := s.uc.CreatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}
	return &subpb.CreatePlanResponse{Message: "Plan created successfully"}, nil
}

func (s *Subscription) UpdatePlan(ctx context.Context, req *subpb.UpdatePlanRequest) (*subpb.UpdatePlanResponse, error) {
	if req.Filter == nil || req.Update == nil {
		return nil, status.Errorf(codes.InvalidArgument, "filter and update fields must be provided")
	}

	update := model.PlanUpdate{
		Name:         optionalString(req.Update.Name),
		Features:     optionalStringSlice(req.Update.Features),
		Price:        optionalFloat(req.Update.Price),
		DurationDays: optionalInt(int(req.Update.DurationDays)),
	}

	err := s.uc.UpdatePlan(ctx, req.Filter.Id, update)
	if err != nil {
		return nil, err
	}
	return &subpb.UpdatePlanResponse{Message: "Plan updated successfully"}, nil
}



func (s *Subscription) GetPlanByID(ctx context.Context, req *subpb.GetPlanByIDRequest) (*subpb.GetPlanByIDResponse, error) {
	plan, err := s.uc.GetPlanByID(ctx, req.Filter.Id)
	if err != nil {
		return nil, err
	}
	return &subpb.GetPlanByIDResponse{Plan: dto.ToProtoPlan(*plan)}, nil
}

func (s *Subscription) GetListOfPlans(ctx context.Context, req *subpb.GetListOfPlansRequest) (*subpb.GetListOfPlansResponse, error) {
	plans, err := s.uc.GetListOfPlans(ctx)
	if err != nil {
		return nil, err
	}
	var protoPlans []*subpb.Plan
	for _, p := range plans {
		protoPlans = append(protoPlans, dto.ToProtoPlan(p))
	}
	return &subpb.GetListOfPlansResponse{Plans: protoPlans}, nil
}


func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalStringSlice(s []string) *[]string {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func optionalFloat(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

func optionalInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
