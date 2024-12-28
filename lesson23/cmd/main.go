package main

import (
	"log"

	"lesson23/handler"
	"lesson23/postgres"
	"lesson23/repository"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	h := handler.NewHandler(orderRepo)

	r := handler.Run(h)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
