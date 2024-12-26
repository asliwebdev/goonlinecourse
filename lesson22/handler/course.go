package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"lesson22/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var course model.Course

	err = json.Unmarshal(body, &course)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if course.Name == "" || course.Tutor == "" {
		http.Error(w, "Course name and tutor are required", http.StatusBadRequest)
		return
	}

	if course.StartedAt.IsZero() {
		course.StartedAt = time.Now()
	}

	courseID := uuid.New().String()
	course.CreatedAt = time.Now()
	course.UpdatedAt = course.CreatedAt

	err = h.courseRepo.CreateCourse(courseID, course)
	if err != nil {
		http.Error(w, "Failed to create course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"course_id": courseID,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("id")
	if courseID == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	course, err := h.courseRepo.GetCourse(courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve course", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var course model.Course
	err = json.Unmarshal(body, &course)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if course.Name == "" || course.Tutor == "" {
		http.Error(w, "Course name and tutor are required", http.StatusBadRequest)
		return
	}

	if course.StartedAt.IsZero() {
		course.StartedAt = time.Now()
	}

	course.UpdatedAt = time.Now()

	courseID := r.URL.Query().Get("id")
	if courseID == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	updatedCourse, err := h.courseRepo.UpdateCourse(courseID, course)
	if err != nil {
		http.Error(w, "Failed to update course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Course updated successfully",
		"course":  updatedCourse,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("id")
	if courseID == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	err := h.courseRepo.DeleteCourse(courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete course", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Course deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetListCourse(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courseRepo.GetListCourse()
	if err != nil {
		http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(courses)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
