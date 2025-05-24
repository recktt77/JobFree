package dto

import (
    "github.com/recktt77/projectProto-definitions/auth_service/genproto/auth"
    "github.com/recktt77/JobFree/internal/model"
    "time"
)

func ToModelUser(req *auth.RegisterUserRequest) *model.User {
    return &model.User{
        Email: req.Email,
        Role:  req.Role,
        Profile: model.UserProfile{
            Name:      req.Profile.Name,
            Bio:       req.Profile.Bio,
            Skills:    req.Profile.Skills,
            AvatarURL: req.Profile.AvatarUrl,
        },
    }
}

func ToProtoProfile(profile model.UserProfile) *auth.UserProfile {
    return &auth.UserProfile{
        Name:      profile.Name,
        Bio:       profile.Bio,
        Skills:    profile.Skills,
        AvatarUrl: profile.AvatarURL,
    }
}

func ToProtoGetProfile(user *model.User) *auth.GetProfileResponse {
    return &auth.GetProfileResponse{
        Id:        user.ID.Hex(),
        Email:     user.Email,
        Role:      user.Role,
        Profile:   ToProtoProfile(user.Profile),
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
    }
}
