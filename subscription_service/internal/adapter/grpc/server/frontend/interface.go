package frontend

import (
	"context"
	subpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
)

type Handler interface {
	Subscribe(ctx context.Context, req *subpb.SubscribeRequest) (*subpb.SubscribeResponse, error)
	CancelSubscription(ctx context.Context, req *subpb.CancelSubscriptionRequest) (*subpb.CancelSubscriptionResponse, error)
	GetSubscriptionStatus(ctx context.Context, req *subpb.GetSubscriptionStatusRequest) (*subpb.GetSubscriptionStatusResponse, error)
	GetSubscription(ctx context.Context, req *subpb.SubscribeRequest) (*subpb.GetSubscriptionResponse, error)

	CreatePlan(ctx context.Context, req *subpb.CreatePlanRequest) (*subpb.CreatePlanResponse, error)
	UpdatePlan(ctx context.Context, req *subpb.UpdatePlanRequest) (*subpb.UpdatePlanResponse, error)
	GetPlanByID(ctx context.Context, req *subpb.GetPlanByIDRequest) (*subpb.GetPlanByIDResponse, error)
	GetListOfPlans(ctx context.Context, req *subpb.GetListOfPlansRequest) (*subpb.GetListOfPlansResponse, error)

}
