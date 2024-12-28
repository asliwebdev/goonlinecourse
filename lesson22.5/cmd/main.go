package main

import (
	"log"

	"lesson22.5/handler"
	"lesson22.5/postgres"
	"lesson22.5/repository"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	h := handler.NewHandler(productRepo, categoryRepo, orderRepo)
	server := handler.Run(h)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
