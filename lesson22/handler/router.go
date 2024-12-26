package handler

import (
	"lesson22/repository"
	"net/http"
)

type Handler struct {
	studentRepo *repository.StudentRepository
	courseRepo  *repository.CourseRepository
	tutorRepo   *repository.TutorRepository
	groupRepo   *repository.GroupRepository
}

func NewHandler(studentRepo *repository.StudentRepository, courseRepo *repository.CourseRepository, tutorRepo *repository.TutorRepository, groupRepo *repository.GroupRepository) *Handler {
	return &Handler{
		studentRepo: studentRepo,
		courseRepo:  courseRepo,
		tutorRepo:   tutorRepo,
		groupRepo:   groupRepo,
	}
}

func Run(handler *Handler) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /student", handler.CreateStudent)
	mux.HandleFunc("GET /student", handler.GetListStudent)
	mux.HandleFunc("GET /student/{id}", handler.GetStudent)
	mux.HandleFunc("PUT /student/{id}", handler.UpdateStudent)
	mux.HandleFunc("DELETE /student/{id}", handler.DeleteStudent)

	mux.HandleFunc("POST /course", handler.CreateCourse)
	mux.HandleFunc("GET /course", handler.GetListCourse)
	mux.HandleFunc("GET /course/{id}", handler.GetCourse)
	mux.HandleFunc("PUT /course/{id}", handler.UpdateCourse)
	mux.HandleFunc("DELETE /course/{id}", handler.DeleteCourse)

	mux.HandleFunc("POST /tutor", handler.CreateTutor)
	mux.HandleFunc("GET /tutor", handler.GetListTutors)
	mux.HandleFunc("GET /tutor/{id}", handler.GetTutor)
	mux.HandleFunc("PUT /tutor/{id}", handler.UpdateTutor)
	mux.HandleFunc("DELETE /tutor/{id}", handler.DeleteTutor)

	mux.HandleFunc("POST /group", handler.CreateGroup)
	mux.HandleFunc("GET /group", handler.GetListGroups)
	mux.HandleFunc("GET /group/{id}", handler.GetGroup)
	mux.HandleFunc("PUT /group/{id}", handler.UpdateGroup)
	mux.HandleFunc("DELETE /group/{id}", handler.DeleteGroup)

	server := &http.Server{Addr: ":8080", Handler: mux}

	return server
}
