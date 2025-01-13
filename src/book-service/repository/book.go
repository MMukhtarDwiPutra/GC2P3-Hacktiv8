package repository

import (
	"book-service/models"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	CreateBook(book models.Book) (*primitive.ObjectID, error)
	GetBookById(id primitive.ObjectID) (*models.Book, error)
	GetAllBooks() (*[]models.Book, error)
	UpdateBook(id primitive.ObjectID, book models.BookRequest) (*mongo.UpdateResult, error)
	DeleteBookById(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type bookRepository struct {
	collection *mongo.Collection
}

// NewBookRepository creates a new instance of bookRepository
func NewBookRepository(collection *mongo.Collection) *bookRepository {
	return &bookRepository{
		collection: collection,
	}
}

func (r *bookRepository) CreateBook(book models.Book) (*primitive.ObjectID, error) {
	// proses insert data ke mongo
	result, err := r.collection.InsertOne(context.Background(), book)
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

func (r *bookRepository) GetBookById(id primitive.ObjectID) (*models.Book, error) {
	var book models.Book

	// Proses get all data
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book) // bson.M -> digunakan untuk filter data based on ...
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *bookRepository) GetAllBooks() (*[]models.Book, error) {
	var books []models.Book

	// Proses get all data
	cursor, err := r.collection.Find(context.Background(), bson.M{}) // bson.M{} -> filter kosong untuk mengambil semua data
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode setiap dokumen ke dalam slice books
	for cursor.Next(context.Background()) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	// Cek apakah cursor menemukan error setelah iterasi
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &books, nil
}

func (r *bookRepository) UpdateBook(id primitive.ObjectID, book models.BookRequest) (*mongo.UpdateResult, error) {
	// Proses Update Data
	result, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": book},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *bookRepository) DeleteBookById(id primitive.ObjectID) (*mongo.DeleteResult, error) {
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
