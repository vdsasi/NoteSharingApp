package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name     string             `bson:"name" json:"name" binding:"required"`
    Email    string             `bson:"email" json:"email" binding:"required,email"`
    Username string             `bson:"username" json:"username" binding:"required"`
    Password string             `bson:"password" json:"-"` // Hide from JSON
}

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    SessionID string `json:"sessionId"`
    User      UserProfileDto `json:"user"`
}

type UserProfileDto struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"oldPassword" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=6"`
}
