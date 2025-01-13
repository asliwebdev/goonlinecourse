package main

import (
	"lesson29/handler"
	"lesson29/postgres"
	"lesson29/repository"
	"lesson29/service"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
	defer redisClient.Close()

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, redisClient)
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/user", userHandler.GetUserById)

	log.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
