package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
    ID             primitive.ObjectID `bson:"_id,omitempty"`
    UserID         primitive.ObjectID `bson:"user_id"`
    SubscriptionID primitive.ObjectID `bson:"subscription_id"`
    Amount         float64            `bson:"amount"`
    Status         string             `bson:"status"` // "pending", "completed", "failed"
    CreatedAt      time.Time          `bson:"created_at"`
    UpdatedAt      time.Time          `bson:"updated_at"`
}

