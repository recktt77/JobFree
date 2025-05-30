package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plan — модель тарифного плана
type Plan struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Features     []string           `bson:"features"`
	Price        float64            `bson:"price"`
	DurationDays int                `bson:"duration_days"`
}

type PlanUpdate struct{
	ID           *primitive.ObjectID
	Name         *string
	Features     *[]string
	Price        *float64
	DurationDays *int
}

type PlanFilter struct{
	ID           *primitive.ObjectID
	Name         *string
	Features     *[]string
	Price        *float64
	DurationDays *int
}
// Subscription — модель подписки пользователя
type Subscription struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	PlanID    primitive.ObjectID `bson:"plan_id"`
	StartDate time.Time          `bson:"start_date"`
	EndDate   time.Time          `bson:"end_date"`
	Active    bool               `bson:"active"`
}


func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func MustObjectIDFromHex(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}
