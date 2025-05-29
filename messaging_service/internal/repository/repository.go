package repository

import (
	"context"
	"messaging_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *model.Message) (primitive.ObjectID, error)
	GetMessagesByConversation(ctx context.Context, conversationID primitive.ObjectID) ([]model.Message, error)
}
