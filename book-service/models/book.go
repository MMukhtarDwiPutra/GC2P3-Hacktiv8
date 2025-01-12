package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	BookID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title         string             `json:"title" gorm:"not null"`
	Author        string             `json:"author" gorm:"not null"`
	PublishedDate time.Time          `json:"published_date" gorm:"not null"`
	Status        string             `json:"status" gorm:"not null"`
	UserID        string             `json:"user_id" gorm:"not null"`
}

type BookRequest struct {
	Title         string    `json:"title" gorm:"not null" validate:"required"`
	Author        string    `json:"author" gorm:"not null" validate:"required"`
	PublishedDate time.Time `json:"published_date" gorm:"not null" validate:"required"`
	Status        string    `json:"status" gorm:"not null" validate:"required"`
	UserID        string    `json:"user_id" gorm:"not null" validate:"required"`
}
