package model

import "context"

type AdminRepository interface {
	BanUser(ctx context.Context, userID, reason string) error
	ApproveProject(ctx context.Context, projectID string) error
	DeleteReview(ctx context.Context, reviewID string) error
	GetStats(ctx context.Context) (*PlatformStats, error)
}
