package frontend

import (
    "github.com/recktt77/JobFree/internal/model"
    "context"
)

type UserUsecase interface {
    Register(ctx context.Context, user *model.User, password string) (string, error)
    Login(ctx context.Context, email, password string) (*model.User, error)
    GetProfile(ctx context.Context, id string) (*model.User, error)
    UpdateProfile(ctx context.Context, id string, profile model.UserProfile) error
}
