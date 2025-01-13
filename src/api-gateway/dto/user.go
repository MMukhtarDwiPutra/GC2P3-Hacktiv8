package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserID   uint32 `json:"_id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password,omitempty" gorm:"not null"`
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
	UserID   primitive.ObjectID `json:"_id" gorm:"primaryKey"`
	Username string             `json:"username" validate:"required"`
}
