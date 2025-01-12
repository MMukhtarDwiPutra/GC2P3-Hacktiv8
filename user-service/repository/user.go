package repository

import (
	"context"
	"log"
	"user-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) (*primitive.ObjectID, error)
}

type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	// Proses get all data
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user) // bson.M -> digunakan untuk filter data based on ...
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user models.User) (*primitive.ObjectID, error) {
	// proses insert data ke mongo
	result, err := r.collection.InsertOne(context.Background(), user)
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
