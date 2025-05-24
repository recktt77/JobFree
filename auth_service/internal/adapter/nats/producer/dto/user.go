package dto

import "github.com/recktt77/JobFree/internal/model"

type UserRegistered struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    CreatedAt string `json:"created_at"`
}

type ProfileUpdated struct {
    ID      string             `json:"id"`
    Profile model.UserProfile  `json:"profile"`
    UpdatedAt string           `json:"updated_at"`
}
