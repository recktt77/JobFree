package dao

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDAO struct {
    Collection *mongo.Collection
}

func NewUserDAO(db *mongo.Database) *UserDAO {
    return &UserDAO{
        Collection: db.Collection("users"),
    }
}

func (dao *UserDAO) Insert(ctx context.Context, user *model.User) (primitive.ObjectID, error) {
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    _, err := dao.Collection.InsertOne(ctx, user)
    return user.ID, err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := dao.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (dao *UserDAO) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
    var user model.User
    err := dao.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (dao *UserDAO) UpdateProfile(ctx context.Context, id primitive.ObjectID, profile model.UserProfile) error {
    _, err := dao.Collection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.M{
            "$set": bson.M{
                "profile":    profile,
                "updated_at": time.Now(),
            },
        },
    )
    return err
}
