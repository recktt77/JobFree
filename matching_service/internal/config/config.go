package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI      string
	Database string
}

type DB struct {
	Conn   *mongo.Database
	Client *mongo.Client
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	log.Printf("Connected to MongoDB: %s", cfg.Database)

	return &DB{
		Conn:   client.Database(cfg.Database),
		Client: client,
	}, nil
}
