package frontend

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/recktt77/JobFree/internal/model"
	"github.com/recktt77/projectProto-definitions/auth_service/genproto/auth"

	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandler struct {
	auth.UnimplementedAuthServiceServer
	service UserUsecase
}

func NewUserHandler(service UserUsecase) *userHandler {
	return &userHandler{service: service}
}

var jwtSecret = []byte("supersecretkey") // лучше взять из .env

func generateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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

	token, err := generateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}

	return &auth.LoginUserResponse{
		Token: token,
	}, nil
}

func (h *userHandler) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	user, err := h.service.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &auth.GetProfileResponse{
		Id:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
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
