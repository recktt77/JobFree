package repository

import (
	"context"
	"messaging_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepo struct {
	messagesColl *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepo{
		messagesColl: db.Collection("messages"),
	}
}

func (r *messageRepo) CreateMessage(ctx context.Context, msg *model.Message) (primitive.ObjectID, error) {
	msg.SentAt = time.Now()
	res, err := r.messagesColl.InsertOne(ctx, msg)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *messageRepo) GetMessagesByConversation(ctx context.Context, conversationID primitive.ObjectID) ([]model.Message, error) {
	cursor, err := r.messagesColl.Find(ctx, bson.M{"conversation_id": conversationID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
