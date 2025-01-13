package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BorrowedBook struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	BookID       string             `json:"book_id" bson:"book_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	BorrowedDate time.Time          `json:"borrowed_date" bson:"borrowed_date"`
	ReturnDate   *time.Time         `json:"return_date" bson:"return_date"`
}

type BorrowedBookRequest struct {
	BookID       string     `json:"book_id" bson:"book_id"`
	UserID       string     `json:"user_id" bson:"user_id"`
	BorrowedDate time.Time  `json:"borrowed_date" bson:"borrowed_date"`
	ReturnDate   *time.Time `json:"return_date" bson:"return_date"`
}

type UpdateBorrowedBook struct {
	ReturnDate *time.Time `json:"return_date" bson:"return_date"`
}
