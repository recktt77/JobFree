package usecase

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/matching_service/internal/models"
	"github.com/recktt77/JobFree/matching_service/internal/repository"
)

type MatchingUseCase interface {
	CreateBid(ctx context.Context, bid *models.Bid) (string, error)
	GetBidsForProject(ctx context.Context, projectID string) ([]*models.Bid, error)
	MatchFreelancers(ctx context.Context, projectID string) ([]*models.Freelancer, error)
}

type MatchingUseCaseImpl struct {
	repo *repository.BidRepository
}

func NewMatchingUseCase(repo *repository.BidRepository) MatchingUseCase {
	return &MatchingUseCaseImpl{repo: repo}
}

func (uc *MatchingUseCaseImpl) CreateBid(ctx context.Context, bid *models.Bid) (string, error) {
	bid.Timestamp = time.Now().Format(time.RFC3339)
	err := uc.repo.Save(ctx, bid)
	if err != nil {
		return "", err
	}
	return bid.BidID, nil
}

func (uc *MatchingUseCaseImpl) GetBidsForProject(ctx context.Context, projectID string) ([]*models.Bid, error) {
	return uc.repo.GetByProjectID(ctx, projectID)
}

func (uc *MatchingUseCaseImpl) MatchFreelancers(ctx context.Context, projectID string) ([]*models.Freelancer, error) {
	// Можно позже реализовать простой алгоритм фильтрации по ключевым словам
	return []*models.Freelancer{}, nil
}
