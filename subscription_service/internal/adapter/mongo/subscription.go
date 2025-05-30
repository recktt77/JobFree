package mongo

import (
	"context"
	"errors"

	"github.com/recktt77/JobFree/subscription_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	planColl         *mongo.Collection
	subscriptionColl *mongo.Collection
}

func NewRepository(db *mongo.Database) *mongoRepo {
	return &mongoRepo{
		planColl:         db.Collection("plans"),
		subscriptionColl: db.Collection("subscriptions"),
	}
}

// ------- PLAN -------
func (r *mongoRepo) CreatePlan(ctx context.Context, plan model.Plan) error {
	_, err := r.planColl.InsertOne(ctx, plan)
	return err
}

func (r *mongoRepo) GetPlanByID(ctx context.Context, filter model.PlanFilter) (model.Plan, error) {
	if filter.ID == nil {
		return model.Plan{}, errors.New("plan ID is required")
	}
	var plan model.Plan
	err := r.planColl.FindOne(ctx, bson.M{"_id": filter.ID}).Decode(&plan)
	return plan, err
}


func (r *mongoRepo) UpdatePlan(ctx context.Context, filter model.PlanFilter, update model.PlanUpdate) error {
	if filter.ID == nil {
		return errors.New("plan ID is required")
	}
	updateDoc := bson.M{}
	if update.Name != nil {
		updateDoc["name"] = *update.Name
	}
	if update.Features != nil {
		updateDoc["features"] = update.Features
	}
	if update.Price != nil {
		updateDoc["price"] = *update.Price
	}
	if update.DurationDays != nil {
		updateDoc["duration_days"] = *update.DurationDays
	}

	_, err := r.planColl.UpdateByID(ctx, filter.ID, bson.M{"$set": updateDoc})
	return err
}


func (r *mongoRepo) GetListOfPlans(ctx context.Context) ([]model.Plan, error) {
	cursor, err := r.planColl.Find(ctx, bson.M{})
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



// ------- SUBSCRIPTION -------
func (r *mongoRepo) CreateSubscription(ctx context.Context, sub model.Subscription) error {
	_, err := r.subscriptionColl.InsertOne(ctx, sub)
	return err
}

func (r *mongoRepo) CancelSubscription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.subscriptionColl.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"active": false}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *mongoRepo) GetUserSubscriptions(ctx context.Context, userID string) ([]model.Subscription, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.subscriptionColl.Find(ctx, bson.M{"user_id": uid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subs []model.Subscription
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *mongoRepo) GetSubscriptionStatus(ctx context.Context, userID string) (*model.Subscription, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var sub model.Subscription
	err = r.subscriptionColl.FindOne(ctx, bson.M{"user_id": uid, "active": true}).Decode(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

