package handler

import (
	"database/sql"
	"encoding/json"
	"lesson22/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) CreateTutor(w http.ResponseWriter, r *http.Request) {
	var tutor model.Tutor
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if tutor.Name == "" || tutor.Subject == "" {
		http.Error(w, "Tutor name and subject are required", http.StatusBadRequest)
		return
	}

	tutor.Id = uuid.New().String()
	tutor.CreatedAt = time.Now()

	err = h.tutorRepo.CreateTutor(&tutor)
	if err != nil {
		http.Error(w, "Failed to create tutor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"tutor_id": tutor.Id,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetTutor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Tutor ID is required", http.StatusBadRequest)
		return
	}

	tutor, err := h.tutorRepo.GetTutor(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tutor not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve tutor", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutor)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateTutor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Tutor ID is required", http.StatusBadRequest)
		return
	}

	var tutor model.Tutor
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if tutor.Name == "" || tutor.Subject == "" {
		http.Error(w, "Tutor name and subject are required", http.StatusBadRequest)
		return
	}

	tutor.Id = id
	updatedTutor, err := h.tutorRepo.UpdateTutor(&tutor)
	if err != nil {
		http.Error(w, "Failed to update tutor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Tutor updated successfully",
		"tutor":   updatedTutor,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteTutor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Tutor ID is required", http.StatusBadRequest)
		return
	}

	err := h.tutorRepo.DeleteTutor(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tutor not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete tutor", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Tutor deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetListTutors(w http.ResponseWriter, r *http.Request) {
	tutors, err := h.tutorRepo.GetListTutors()
	if err != nil {
		http.Error(w, "Failed to retrieve tutors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutors)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
