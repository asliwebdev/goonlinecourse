package main

import (
	"lesson31/models"
	"lesson31/repository"
	"lesson31/repository/mongodb" //
	"log"
)

func main() {
	client, err := mongodb.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	dbName := "go_online"       
	collectionName := "users"         
	userRepo := repository.NewUserRepo(client, dbName, collectionName)

	user := models.User{
		Name:  "Jane Smith",
		Role:  "user",
		Email: "janesmith@example.com",
	}

	err = userRepo.CreateUser(user)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
}