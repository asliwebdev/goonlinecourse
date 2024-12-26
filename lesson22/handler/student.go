package handler

import (
	"encoding/json"
	"io"
	"lesson22/model"
	"log"
	"net/http"
)

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	student := model.CreateStudent{}

	err = json.Unmarshal(body, &student)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.studentRepo.CreateStudent(&student)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) GetListStudent(w http.ResponseWriter, r *http.Request) {
	students, err := h.studentRepo.GetListStudent()
	if err != nil {
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resp, err := h.studentRepo.GetStudent(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var student model.Student
	err = json.Unmarshal(body, &student)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if student.Name == "" || student.LastName == "" {
		http.Error(w, "Student name and last name are required", http.StatusBadRequest)
		return
	}

	studentID := r.URL.Query().Get("id")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	updatedStudent, err := h.studentRepo.UpdateStudent(student)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Student updated successfully",
		"student": updatedStudent,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("id")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	err := h.studentRepo.DeleteStudent(studentID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Student deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}
