package repository

import (
	"context"
	"errors"

	"github.com/recktt77/JobFree/matching_service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidRepository struct {
	collection *mongo.Collection
}

func NewBidRepository(db *mongo.Database) *BidRepository {
	return &BidRepository{
		collection: db.Collection("bids"),
	}
}

func (r *BidRepository) Save(ctx context.Context, bid *models.Bid) error {
	_, err := r.collection.InsertOne(ctx, bid)
	return err
}

func (r *BidRepository) GetByProjectID(ctx context.Context, projectID string) ([]*models.Bid, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"projectid": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bids []*models.Bid
	for cursor.Next(ctx) {
		var bid models.Bid
		if err := cursor.Decode(&bid); err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	if len(bids) == 0 {
		return nil, errors.New("no bids found")
	}

	return bids, nil
}
