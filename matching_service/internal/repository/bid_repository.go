package repository

import (
	"context"

	"github.com/recktt77/JobFree/matching_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidRepository interface {
	Create(ctx context.Context, bid *model.Bid) error
	GetByProjectID(ctx context.Context, projectID string) ([]*model.Bid, error)
}

type bidRepository struct {
	collection *mongo.Collection
}

func NewBidRepository(db *mongo.Database) BidRepository {
	return &bidRepository{
		collection: db.Collection("bids"),
	}
}

func (r *bidRepository) Create(ctx context.Context, bid *model.Bid) error {
	_, err := r.collection.InsertOne(ctx, bid)
	return err
}

func (r *bidRepository) GetByProjectID(ctx context.Context, projectID string) ([]*model.Bid, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"projectid": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bids []*model.Bid
	for cursor.Next(ctx) {
		var bid model.Bid
		if err := cursor.Decode(&bid); err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}
