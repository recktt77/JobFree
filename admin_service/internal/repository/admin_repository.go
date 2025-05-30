package repository

import (
	"context"

	"github.com/recktt77/JobFree/admin_service/internal/model"
)

type AdminRepository interface {
	BanUser(ctx context.Context, userID, reason string) error
	DeleteReview(ctx context.Context, reviewID, moderatorID string) error
	ModerateProject(ctx context.Context, projectID, action string) error
	GetPlatformStats(ctx context.Context) (*model.PlatformStats, error)
}
