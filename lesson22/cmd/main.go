package main

import (
	"lesson22/handler"
	"lesson22/postgres"
	"lesson22/repository"
	"log"
)

func main() {
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	tutorRepo := repository.NewTutorRepository(db)
	groupRepo := repository.NewGroupRepository(db)

	h := handler.NewHandler(studentRepo, courseRepo, tutorRepo, groupRepo)

	mux := handler.Run(h)

	err = mux.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
