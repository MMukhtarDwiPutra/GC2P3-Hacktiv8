package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookRegister struct {
	Title         string `json:"title" gorm:"not null" validate:"required"`
	Author        string `json:"author" gorm:"not null" validate:"required"`
	PublishedDate string `json:"published_date" gorm:"not null" validate:"required"`
	Status        string `json:"status" gorm:"not null" validate:"required"`
}

type BookResponse struct {
	BookID        primitive.ObjectID `json:"_id" gorm:"primaryKey"`
	Title         string             `json:"title" gorm:"not null"`
	Author        string             `json:"author" gorm:"not null"`
	PublishedDate string             `json:"published_date" gorm:"not null"`
	Status        string             `json:"status" gorm:"not null"`
	UserID        string             `json:"user_id" gorm:"not null"`
}

type BorrowResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	BookID       string             `json:"book_id" bson:"book_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	BorrowedDate string             `json:"borrowed_date" bson:"borrowed_date"`
	ReturnDate   string             `json:"return_date" bson:"return_date"`
}

// ErrorResponse represents a standardized error response structure.
type SuccessResponse struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
	Data    any    `json:"details"` // Optional details for the error
}

type RegisterResponse struct {
	UserID   primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	Username string             `json:"username" validate:"required"`
}
