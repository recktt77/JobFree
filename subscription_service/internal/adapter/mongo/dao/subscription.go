package dao

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/recktt77/JobFree/subscription_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionDAO struct {
	subscriptions *mongo.Collection
	plans         *mongo.Collection
}

func NewSubscriptionDAO(db *mongo.Database) *SubscriptionDAO {
	return &SubscriptionDAO{
		subscriptions: db.Collection("subscriptions"),
		plans:         db.Collection("plans"),
	}
}

func (d *SubscriptionDAO) Subscribe(ctx context.Context, userID, planID string) (*model.Subscription, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	pid, err := primitive.ObjectIDFromHex(planID)
	if err != nil {
		return nil, err
	}

	var plan model.Plan
	if err := d.plans.FindOne(ctx, bson.M{"_id": pid}).Decode(&plan); err != nil {
		return nil, err
	}

	sub := model.Subscription{
		ID:        primitive.NewObjectID(),
		UserID:    uid,
		PlanID:    pid,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, plan.DurationDays),
		Active:    true,
	}

	_, err = d.subscriptions.InsertOne(ctx, sub)
	if err != nil {
		return nil, err
	}

	log.Printf("Created subscription with id: %s", sub.ID.Hex())
	return &sub, nil
}

func (d *SubscriptionDAO) CancelSubscription(ctx context.Context, subscriptionID string) error {
	id, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"active": false}}

	res, err := d.subscriptions.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (d *SubscriptionDAO) GetSubscription(ctx context.Context, userID string) ([]model.Subscription, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	filter := bson.M{"user_id": uid}
	cursor, err := d.subscriptions.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []model.Subscription
	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}


func (d *SubscriptionDAO) GetSubscriptionStatus(ctx context.Context, userID string) (*model.Subscription, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var sub model.Subscription
	err = d.subscriptions.FindOne(ctx, bson.M{
		"user_id": uid,
		"active":  true,
	}).Decode(&sub)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (d *SubscriptionDAO) GetSubscriptionByID(ctx context.Context, id string) (model.Subscription, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("invalid ObjectID: %w", err)
	}

	var sub model.Subscription
	err = d.subscriptions.FindOne(ctx, bson.M{"_id": objectID}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Subscription{}, errors.New("subscription not found")
		}
		return model.Subscription{}, err
	}
	return sub, nil
}

func (d *SubscriptionDAO) GetListOfPlans(ctx context.Context) ([]model.Plan, error) {
	cursor, err := d.plans.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []model.Plan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

