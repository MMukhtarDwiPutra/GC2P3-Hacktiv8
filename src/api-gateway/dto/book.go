package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	BookID        primitive.ObjectID `json:"_id" gorm:"primaryKey"`
	Title         string             `json:"title" gorm:"not null"`
	Author        string             `json:"author" gorm:"not null"`
	PublishedDate string             `json:"published_date" gorm:"not null"`
	Status        string             `json:"status" gorm:"not null"`
	UserID        string             `json:"user_id" gorm:"not null"`
}

type BookRequest struct {
	Title         string `json:"title" gorm:"not null" validate:"required"`
	Author        string `json:"author" gorm:"not null" validate:"required"`
	PublishedDate string `json:"published_date" gorm:"not null" validate:"required"`
	Status        string `json:"status" gorm:"not null" validate:"required"`
}
