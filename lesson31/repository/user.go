package repository

import (
	"context"
	"fmt"
	"log"

	"lesson31/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Collection
}

func NewUserRepo(client *mongo.Client, dbName, collectionName string) *UserRepo {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepo{DB: collection}
}

func (u *UserRepo) CreateUser(user models.User) error {
	id := uuid.NewString()

	userDoc := bson.M{
		"id":    id,
		"name":  user.Name,
		"role":  user.Role,
		"email": user.Email,
	}
	
	_, err := u.DB.InsertOne(context.Background(), userDoc)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}

	fmt.Printf("User with id %s created successfully.\n", id)
	return nil
}
