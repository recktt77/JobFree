package frontend

import (
    "github.com/recktt77/projectProto-definitions/auth_service/genproto/auth"
    "github.com/recktt77/JobFree/internal/model"
    "context"
    "time"

    _"go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandler struct {
    auth.UnimplementedAuthServiceServer
    service UserUsecase
}

func NewUserHandler(service UserUsecase) *userHandler {
    return &userHandler{service: service}
}

func (h *userHandler) RegisterUser(ctx context.Context, req *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
    user := &model.User{
        Email: req.Email,
        Role:  req.Role,
        Profile: model.UserProfile{
            Name:      req.Profile.Name,
            Bio:       req.Profile.Bio,
            Skills:    req.Profile.Skills,
            AvatarURL: req.Profile.AvatarUrl,
        },
    }

    id, err := h.service.Register(ctx, user, req.Password)
    if err != nil {
        return nil, err
    }

    return &auth.RegisterUserResponse{
        Id:      id,
        Message: "User registered successfully",
    }, nil
}

func (h *userHandler) LoginUser(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
    user, err := h.service.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }

    return &auth.LoginUserResponse{
        Token: user.ID.Hex(),
    }, nil
}

func (h *userHandler) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
    user, err := h.service.GetProfile(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return &auth.GetProfileResponse{
        Id:        user.ID.Hex(),
        Email:     user.Email,
        Role:      user.Role,
        Profile: &auth.UserProfile{
            Name:      user.Profile.Name,
            Bio:       user.Profile.Bio,
            Skills:    user.Profile.Skills,
            AvatarUrl: user.Profile.AvatarURL,
        },
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
    }, nil
}

func (h *userHandler) UpdateProfile(ctx context.Context, req *auth.UpdateProfileRequest) (*auth.UpdateProfileResponse, error) {
    profile := model.UserProfile{
        Name:      req.Profile.Name,
        Bio:       req.Profile.Bio,
        Skills:    req.Profile.Skills,
        AvatarURL: req.Profile.AvatarUrl,
    }

    err := h.service.UpdateProfile(ctx, req.Id, profile)
    if err != nil {
        return nil, err
    }

    return &auth.UpdateProfileResponse{Message: "Profile updated"}, nil
}
