package mongo

import (
    "github.com/recktt77/JobFree/internal/adapter/mongo/dao"
    "github.com/recktt77/JobFree/internal/model"
    "context"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
    dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
    return &UserRepository{dao: dao}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (primitive.ObjectID, error) {
    return r.dao.Insert(ctx, user)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    return r.dao.FindByEmail(ctx, email)
}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
    return r.dao.FindByID(ctx, id)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id primitive.ObjectID, profile model.UserProfile) error {
    return r.dao.UpdateProfile(ctx, id, profile)
}
