package usecase

import (
    "github.com/recktt77/JobFree/internal/model"
    "context"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) (primitive.ObjectID, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
    UpdateProfile(ctx context.Context, id primitive.ObjectID, profile model.UserProfile) error
}

type Cache interface {
    GetUser(ctx context.Context, id string) (*model.User, error)
    SetUser(ctx context.Context, user *model.User) error
}

type EventPublisher interface {
    PublishUserRegistered(user *model.User) error
    PublishProfileUpdated(userID string, profile model.UserProfile) error
}
