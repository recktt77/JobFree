package dto

import (
	"github.com/recktt77/JobFree/subscription_service/internal/model"
	eventpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
)


func ToPlanPayload(plan model.Plan) *eventpb.PlanPayload {
	return &eventpb.PlanPayload{
		Id:           plan.ID.Hex(),
		Name:         plan.Name,
		Features:     plan.Features,
		Price:        plan.Price,
		DurationDays: int32(plan.DurationDays),
	}
}