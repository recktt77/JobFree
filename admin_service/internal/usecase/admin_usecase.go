package usecase

import (
	"context"

	"github.com/recktt77/JobFree/admin_service/internal/events"
	"github.com/recktt77/JobFree/admin_service/internal/metrics"
	"github.com/recktt77/JobFree/admin_service/internal/model"
	"github.com/recktt77/JobFree/admin_service/internal/repository"
)

type AdminUseCase interface {
	BanUser(ctx context.Context, userID, reason string) error
	DeleteReview(ctx context.Context, reviewID, moderatorID string) error
	ModerateProject(ctx context.Context, projectID, action string) error
	GetPlatformStats(ctx context.Context) (*model.PlatformStats, error)
}

type adminUseCase struct {
	repo      repository.AdminRepository
	publisher *events.EventPublisher
}

func NewAdminUseCase(r repository.AdminRepository, publisher *events.EventPublisher) AdminUseCase {
	return &adminUseCase{
		repo:      r,
		publisher: publisher,
	}
}

func (uc *adminUseCase) BanUser(ctx context.Context, userID, reason string) error {
	// 1. Баним в БД
	err := uc.repo.BanUser(ctx, userID, reason)
	if err != nil {
		return err
	}

	// 2. Увеличиваем метрику
	metrics.BanUserCounter.Inc()

	// 3. Публикуем событие в NATS
	go uc.publisher.PublishUserBanned(ctx, userID, reason)

	return nil
}

func (uc *adminUseCase) DeleteReview(ctx context.Context, reviewID, moderatorID string) error {
	return uc.repo.DeleteReview(ctx, reviewID, moderatorID)
}

func (uc *adminUseCase) ModerateProject(ctx context.Context, projectID, action string) error {
	return uc.repo.ModerateProject(ctx, projectID, action)
}

func (uc *adminUseCase) GetPlatformStats(ctx context.Context) (*model.PlatformStats, error) {
	return uc.repo.GetPlatformStats(ctx)
}
