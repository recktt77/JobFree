package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ConversationID primitive.ObjectID `bson:"conversation_id"`
	SenderID       primitive.ObjectID `bson:"sender_id"`
	Content        string             `bson:"content"`
	SentAt         time.Time          `bson:"sent_at"`
}

type Conversation struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	UserIDs   []primitive.ObjectID `bson:"user_ids"`
	CreatedAt time.Time            `bson:"created_at"`
}

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Type      string             `bson:"type"`
	Message   string             `bson:"message"`
	Read      bool               `bson:"read"`
	CreatedAt time.Time          `bson:"created_at"`
}
