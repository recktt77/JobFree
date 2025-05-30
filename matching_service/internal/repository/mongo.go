package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Bids *mongo.Collection
}

func NewMongoRepository(uri string) *MongoRepository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return &MongoRepository{
		Bids: client.Database("jobfree-matching").Collection("bids"),
	}
}
