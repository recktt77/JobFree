package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    Email        string             `bson:"email"`
    PasswordHash string             `bson:"password_hash"`
    Role         string             `bson:"role"`
    Profile      UserProfile        `bson:"profile"`
    CreatedAt    time.Time          `bson:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at"`
}

type UserProfile struct {
    Name      string   `bson:"name"`
    Bio       string   `bson:"bio"`
    Skills    []string `bson:"skills"`
    AvatarURL string   `bson:"avatar_url"`
}
