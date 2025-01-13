package repository

import (
	"borrow-service/models"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BorrowRepository interface {
	CreateBorrow(borrow models.BorrowedBook) (*primitive.ObjectID, error)
	GetBorrowById(id primitive.ObjectID) (*models.BorrowedBook, error)
	GetAllBorrowsByUser(id primitive.ObjectID) (*[]models.BorrowedBook, error)
	UpdateBorrow(id primitive.ObjectID, borrow models.BorrowedBookRequest) (*mongo.UpdateResult, error)
	DeleteBorrowById(id primitive.ObjectID) (*mongo.DeleteResult, error)
	UpdateReturnBorrow(id primitive.ObjectID, borrow models.UpdateBorrowedBook) (*mongo.UpdateResult, error)
}

type borrowRepository struct {
	collection *mongo.Collection
}

// NewBorrowRepository creates a new instance of borrowRepository
func NewBorrowRepository(collection *mongo.Collection) *borrowRepository {
	return &borrowRepository{
		collection: collection,
	}
}

func (r *borrowRepository) CreateBorrow(borrow models.BorrowedBook) (*primitive.ObjectID, error) {
	// proses insert data ke mongo
	result, err := r.collection.InsertOne(context.Background(), borrow)
	if err != nil {
		return nil, err
	}

	// Extracting the Inserted ID as ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("Failed to assert InsertedID as ObjectID")
	}

	return &insertedID, nil
}

func (r *borrowRepository) GetBorrowById(id primitive.ObjectID) (*models.BorrowedBook, error) {
	var borrow models.BorrowedBook

	// Proses get all data
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&borrow) // bson.M -> digunakan untuk filter data based on ...
	if err != nil {
		return nil, err
	}

	return &borrow, nil
}

func (r *borrowRepository) GetAllBorrowsByUser(id primitive.ObjectID) (*[]models.BorrowedBook, error) {
	// Convert ObjectID to string
	userId := id.Hex()
	var borrows []models.BorrowedBook

	// Use Find for retrieving multiple records
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background()) // Ensure the cursor is closed to prevent resource leaks

	// Decode all documents into the slice
	if err = cursor.All(context.Background(), &borrows); err != nil {
		return nil, err
	}

	return &borrows, nil
}

func (r *borrowRepository) GetAllBorrows() (*[]models.BorrowedBook, error) {
	var borrows []models.BorrowedBook

	// Proses get all data
	cursor, err := r.collection.Find(context.Background(), bson.M{}) // bson.M{} -> filter kosong untuk mengambil semua data
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode setiap dokumen ke dalam slice borrows
	for cursor.Next(context.Background()) {
		var borrow models.BorrowedBook
		if err := cursor.Decode(&borrow); err != nil {
			return nil, err
		}
		borrows = append(borrows, borrow)
	}

	// Cek apakah cursor menemukan error setelah iterasi
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &borrows, nil
}

func (r *borrowRepository) UpdateBorrow(id primitive.ObjectID, borrow models.BorrowedBookRequest) (*mongo.UpdateResult, error) {
	// Proses Update Data
	result, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": borrow},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *borrowRepository) UpdateReturnBorrow(id primitive.ObjectID, borrow models.UpdateBorrowedBook) (*mongo.UpdateResult, error) {
	// Proses Update Data
	result, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": borrow},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *borrowRepository) DeleteBorrowById(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	// Proses Delete Data
	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Untuk return 404
	if result.DeletedCount == 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Data Not Found")
	}

	return result, nil
}
