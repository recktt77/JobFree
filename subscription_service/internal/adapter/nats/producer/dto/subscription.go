package dto

import (
	"github.com/recktt77/JobFree/subscription_service/internal/model"
	eventpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func ToSubscriptionPayload(sub model.Subscription) *eventpb.SubscriptionPayload {
	return &eventpb.SubscriptionPayload{
		Id:        sub.ID.Hex(),
		UserId:    sub.UserID.Hex(),
		PlanId:    sub.PlanID.Hex(),
		StartDate: timestamppb.New(sub.StartDate),
		EndDate:   timestamppb.New(sub.EndDate),
		Active:    sub.Active,
	}
}
