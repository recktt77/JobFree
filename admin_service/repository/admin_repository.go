package repository

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/admin_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminRepository interface {
	BanUser(ctx context.Context, userID string, reason string) error
	DeleteReview(ctx context.Context, reviewID string) error
	ModerateProject(ctx context.Context, projectID string, action string) error
	GetStats(ctx context.Context) (*model.PlatformStats, error)
}

type adminRepo struct {
	db *MongoDB
}

func NewAdminRepository(db *MongoDB) AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) BanUser(ctx context.Context, userID string, reason string) error {
	collection := r.db.Database.Collection("banned_users")
	_, err := collection.InsertOne(ctx, bson.M{
		"user_id":   userID,
		"reason":    reason,
		"banned_at": time.Now(),
	})
	return err
}

func (r *adminRepo) DeleteReview(ctx context.Context, reviewID string) error {
	collection := r.db.Database.Collection("reviews")
	_, err := collection.DeleteOne(ctx, bson.M{"review_id": reviewID})
	return err
}

func (r *adminRepo) ModerateProject(ctx context.Context, projectID string, action string) error {
	collection := r.db.Database.Collection("projects")
	update := bson.M{"$set": bson.M{"status": action}}
	_, err := collection.UpdateOne(ctx, bson.M{"project_id": projectID}, update)
	return err
}

func (r *adminRepo) GetStats(ctx context.Context) (*model.PlatformStats, error) {
	users, _ := r.db.Database.Collection("users").CountDocuments(ctx, bson.M{})
	banned, _ := r.db.Database.Collection("banned_users").CountDocuments(ctx, bson.M{})
	projects, _ := r.db.Database.Collection("projects").CountDocuments(ctx, bson.M{"status": "active"})
	reviews, _ := r.db.Database.Collection("reviews").CountDocuments(ctx, bson.M{})

	return &model.PlatformStats{
		TotalUsers:     int32(users),
		BannedUsers:    int32(banned),
		ActiveProjects: int32(projects),
		TotalReviews:   int32(reviews),
	}, nil
}
