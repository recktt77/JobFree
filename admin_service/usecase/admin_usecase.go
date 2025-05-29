package usecase

import (
	"context"

	"github.com/recktt77/JobFree/admin_service/internal/model"
)

type AdminUsecase struct {
	Repo model.AdminRepository
}

func NewAdminUsecase(repo model.AdminRepository) *AdminUsecase {
	return &AdminUsecase{Repo: repo}
}

func (u *AdminUsecase) BanUser(ctx context.Context, userID, reason string) error {
	return u.Repo.BanUser(ctx, userID, reason)
}

func (u *AdminUsecase) ApproveOrRejectProject(ctx context.Context, projectID, action string) error {
	return u.Repo.ApproveOrRejectProject(ctx, projectID, action)
}

func (u *AdminUsecase) DeleteReview(ctx context.Context, reviewID, moderatorID string) error {
	return u.Repo.DeleteReview(ctx, reviewID, moderatorID)
}

func (u *AdminUsecase) GetStats(ctx context.Context) (*model.PlatformStats, error) {
	return u.Repo.GetPlatformStats(ctx)
}
