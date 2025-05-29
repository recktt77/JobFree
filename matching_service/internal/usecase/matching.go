package usecase

import (
	"context"
	"time"

	"github.com/recktt77/JobFree/matching_service/internal/cache"
	"github.com/recktt77/JobFree/matching_service/internal/events"
	"github.com/recktt77/JobFree/matching_service/internal/model"
	"github.com/recktt77/JobFree/matching_service/internal/repository"
)

type MatchingUseCase struct {
	BidRepo   repository.BidRepository
	Redis     *cache.RedisCache
	Publisher *events.Publisher
}

func NewMatchingUseCase(bidRepo repository.BidRepository, redis *cache.RedisCache, publisher *events.Publisher) MatchingUseCase {
	return MatchingUseCase{
		BidRepo:   bidRepo,
		Redis:     redis,
		Publisher: publisher,
	}
}

// 💾 Создаёт заявку и сохраняет в MongoDB
func (uc MatchingUseCase) CreateBid(ctx context.Context, bid *model.Bid) (string, error) {
	bid.Timestamp = time.Now().Format(time.RFC3339)
	err := uc.Publisher.PublishBidCreated(ctx, bid)

	if err != nil {
		return "", err
	}
	return bid.BidID, nil
}

// 📤 Получает все заявки по project_id
func (uc MatchingUseCase) GetBidsForProject(ctx context.Context, projectID string) ([]model.Bid, error) {
	bidPtrs, err := uc.BidRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Преобразуем []*model.Bid в []model.Bid
	bids := make([]model.Bid, 0, len(bidPtrs))
	for _, ptr := range bidPtrs {
		if ptr != nil {
			bids = append(bids, *ptr)
		}
	}

	return bids, nil
}

// 🎯 Фильтрация по скиллам через Redis
func (uc MatchingUseCase) MatchFreelancers(ctx context.Context, requiredSkills []string) ([]model.Freelancer, error) {
	allFreelancers, err := uc.Redis.GetAllFreelancers()
	if err != nil {
		return nil, err
	}

	var matched []model.Freelancer
	for _, f := range allFreelancers {
		rawSkills, ok := f["skills"].([]interface{})
		if !ok {
			continue
		}
		var skills []string
		for _, s := range rawSkills {
			skills = append(skills, s.(string))
		}
		if hasSkillOverlap(skills, requiredSkills) {
			matched = append(matched, model.Freelancer{
				FreelancerID: f["freelancer_id"].(string),
				Name:         f["name"].(string),
				Skills:       skills,
			})
		}
	}

	return matched, nil
}

func hasSkillOverlap(userSkills, requiredSkills []string) bool {
	skillMap := make(map[string]bool)
	for _, s := range userSkills {
		skillMap[s] = true
	}
	for _, r := range requiredSkills {
		if skillMap[r] {
			return true
		}
	}
	return false
}
