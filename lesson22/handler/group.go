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

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var group model.Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if group.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	groupID := uuid.New().String()
	group.Id = groupID
	group.CreatedAt = time.Now()

	err = h.groupRepo.CreateGroup(group)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"group_id": groupID,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	group, err := h.groupRepo.GetGroup(groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve group", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var group model.Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if group.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	group.Id = groupID
	updatedGroup, err := h.groupRepo.UpdateGroup(group)
	if err != nil {
		http.Error(w, "Failed to update group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Group updated successfully",
		"group":   updatedGroup,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	err := h.groupRepo.DeleteGroup(groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete group", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Group deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupRepo.GetListGroups()
	if err != nil {
		http.Error(w, "Failed to retrieve groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
