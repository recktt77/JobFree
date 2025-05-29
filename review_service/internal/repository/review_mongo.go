package repository

import (
	"context"
	"time"

	"review_service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	Create(ctx context.Context, r *model.Review) (primitive.ObjectID, error)
	GetByProjectID(ctx context.Context, id primitive.ObjectID) ([]model.Review, error)
	GetByRevieweeID(ctx context.Context, id primitive.ObjectID) ([]model.Review, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Review, error)
	Moderate(ctx context.Context, reviewID primitive.ObjectID, action string) error
}

type reviewRepo struct {
	col *mongo.Collection
}

func NewReviewRepo(db *mongo.Database) ReviewRepository {
	return &reviewRepo{col: db.Collection("reviews")}
}

func (r *reviewRepo) Create(ctx context.Context, rev *model.Review) (primitive.ObjectID, error) {
	rev.CreatedAt = time.Now()
	res, err := r.col.InsertOne(ctx, rev)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *reviewRepo) GetByProjectID(ctx context.Context, id primitive.ObjectID) ([]model.Review, error) {
	var reviews []model.Review
	cursor, err := r.col.Find(ctx, bson.M{"project_id": id, "hidden": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepo) GetByRevieweeID(ctx context.Context, id primitive.ObjectID) ([]model.Review, error) {
	var reviews []model.Review
	cursor, err := r.col.Find(ctx, bson.M{"reviewee_id": id, "hidden": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Review, error) {
	var review model.Review
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) Moderate(ctx context.Context, reviewID primitive.ObjectID, action string) error {
	switch action {
	case "delete":
		_, err := r.col.DeleteOne(ctx, bson.M{"_id": reviewID})
		return err
	case "hide":
		_, err := r.col.UpdateOne(ctx, bson.M{"_id": reviewID}, bson.M{"$set": bson.M{"hidden": true}})
		return err
	case "approve":
		_, err := r.col.UpdateOne(ctx, bson.M{"_id": reviewID}, bson.M{"$set": bson.M{"hidden": false}})
		return err
	default:
		return nil
	}
}
