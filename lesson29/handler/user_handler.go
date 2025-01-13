package handler

import (
	"encoding/json"
	"lesson29/service"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("id")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, source, err := h.userService.GetUserById(r.Context(), userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"source": source,
		"data":   user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
