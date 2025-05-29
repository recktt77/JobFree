package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/recktt77/JobFree/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo   UserRepository
	cache  Cache
	events EventPublisher
}

func NewUserUseCase(repo UserRepository, cache Cache, events EventPublisher) *UserUseCase {
	if repo == nil || cache == nil || events == nil {
		panic("NewUserUseCase: repo/cache/events is nil") // временно для отладки
	}

	return &UserUseCase{
		repo:   repo,
		cache:  cache,
		events: events,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, user *model.User, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user.PasswordHash = string(hash)

	id, err := uc.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	// 🔥 Добавляем вызов NATS-публикации
	go func() {
		if err := uc.events.PublishUserRegistered(user); err != nil {
			log.Printf("failed to publish user.registered: %v", err)
		}
	}()

	return id.Hex(), nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	_ = uc.cache.SetUser(ctx, user)
	return user, nil
}

func (uc *UserUseCase) GetProfile(ctx context.Context, id string) (*model.User, error) {
	if cached, err := uc.cache.GetUser(ctx, id); err == nil {
		return cached, nil
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := uc.repo.GetByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	_ = uc.cache.SetUser(ctx, user)
	return user, nil
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, id string, profile model.UserProfile) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := uc.repo.UpdateProfile(ctx, objID, profile); err != nil {
		return err
	}

	user, err := uc.repo.GetByID(ctx, objID)
	if err == nil {
		user.Profile = profile
		user.UpdatedAt = time.Now()
		_ = uc.cache.SetUser(ctx, user)
	}

	return uc.events.PublishProfileUpdated(id, profile)
}
