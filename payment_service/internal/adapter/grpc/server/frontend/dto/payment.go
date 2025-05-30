package dto

import (
	"github.com/recktt77/JobFree/payment_service/internal/model"
	pb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToProto(p *model.Payment) *pb.Payment {
	return &pb.Payment{
		Id:             p.ID.Hex(),
		UserId:         p.UserID.Hex(),
		SubscriptionId: p.SubscriptionID.Hex(),
		Amount:         p.Amount,
		Status:         p.Status,
		CreatedAt:      p.CreatedAt.String(),
		UpdatedAt:      p.UpdatedAt.String(),
	}
}

func FromProto(p *pb.Payment) (*model.Payment, error) {
	id, err := primitive.ObjectIDFromHex(p.Id)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(p.UserId)
	if err != nil {
		return nil, err
	}
	subID, err := primitive.ObjectIDFromHex(p.SubscriptionId)
	if err != nil {
		return nil, err
	}

	return &model.Payment{
		ID:             id,
		UserID:         userID,
		SubscriptionID: subID,
		Amount:         p.Amount,
		Status:         p.Status,
	}, nil
}
