package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID  primitive.ObjectID `bson:"project_id"`
	ReviewerID primitive.ObjectID `bson:"reviewer_id"`
	RevieweeID primitive.ObjectID `bson:"reviewee_id"`
	Rating     int32              `bson:"rating"`
	Comment    string             `bson:"comment"`
	CreatedAt  time.Time          `bson:"created_at"`
	Hidden     bool               `bson:"hidden"`
}
