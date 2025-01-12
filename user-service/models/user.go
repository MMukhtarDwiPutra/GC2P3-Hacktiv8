package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserID   primitive.ObjectID `json:"_id"`
	Username string             `json:"username"`
	Password string             `json:"password,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UserID   primitive.ObjectID `json:"_id"`
	Username string             `json:"username" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
