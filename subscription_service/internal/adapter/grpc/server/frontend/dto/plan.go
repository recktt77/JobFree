package dto

import (
	"github.com/recktt77/JobFree/subscription_service/internal/model"
	subpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
)

func FromProtoToPlan(req *subpb.CreatePlanRequest) model.Plan {
	return model.Plan{
		Name:         req.Plan.GetName(),
		Features:     req.Plan.GetFeatures(),
		Price:        req.Plan.GetPrice(),
		DurationDays: int(req.Plan.GetDurationDays()),
	}
}

func ToProtoPlan(plan model.Plan) *subpb.Plan {
	return &subpb.Plan{
		Id:           plan.ID.Hex(),
		Name:         plan.Name,
		Features:     plan.Features,
		Price:        plan.Price,
		DurationDays: int32(plan.DurationDays),
	}
}
