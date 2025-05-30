package repository

import (
	"context"
	"errors"

	"github.com/recktt77/JobFree/admin_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	db *mongo.Database
}

func NewAdminRepository(db *mongo.Database) AdminRepository {
	return &mongoRepo{db: db}
}

func (r *mongoRepo) BanUser(ctx context.Context, userID, reason string) error {
	_, err := r.db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"banned": true, "ban_reason": reason}},
	)
	return err
}

func (r *mongoRepo) DeleteReview(ctx context.Context, reviewID, moderatorID string) error {
	_, err := r.db.Collection("reviews").UpdateOne(ctx,
		bson.M{"_id": reviewID},
		bson.M{"$set": bson.M{"deleted": true, "moderator_id": moderatorID}},
	)
	return err
}

func (r *mongoRepo) ModerateProject(ctx context.Context, projectID, action string) error {
	if action != "approve" && action != "reject" {
		return errors.New("invalid action")
	}
	_, err := r.db.Collection("projects").UpdateOne(ctx,
		bson.M{"_id": projectID},
		bson.M{"$set": bson.M{"status": action}},
	)
	return err
}

func (r *mongoRepo) GetPlatformStats(ctx context.Context) (*model.PlatformStats, error) {
	usersCount, _ := r.db.Collection("users").CountDocuments(ctx, bson.M{})
	bannedCount, _ := r.db.Collection("users").CountDocuments(ctx, bson.M{"banned": true})
	projectsCount, _ := r.db.Collection("projects").CountDocuments(ctx, bson.M{"status": "approved"})
	reviewsCount, _ := r.db.Collection("reviews").CountDocuments(ctx, bson.M{})
	return &model.PlatformStats{
		TotalUsers:     int32(usersCount),
		BannedUsers:    int32(bannedCount),
		ActiveProjects: int32(projectsCount),
		TotalReviews:   int32(reviewsCount),
	}, nil
}
