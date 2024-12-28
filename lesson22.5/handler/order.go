package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"lesson22.5/model"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var order model.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if order.CustomerID == "" {
		http.Error(w, "CustomerID is required", http.StatusBadRequest)
		return
	}
	if len(order.Products) == 0 {
		http.Error(w, "At least one product is required", http.StatusBadRequest)
		return
	}
	if order.TotalAmount <= 0 {
		http.Error(w, "TotalAmount must be greater than zero", http.StatusBadRequest)
		return
	}
	if order.ShippingAddress == "" {
		http.Error(w, "ShippingAddress is required", http.StatusBadRequest)
		return
	}
	if order.PaymentMethod == "" {
		http.Error(w, "PaymentMethod is required", http.StatusBadRequest)
		return
	}

	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	err = h.orderRepo.CreateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Order created successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderRepo.GetOrderById(id)
	if err != nil {
		http.Error(w, "Failed to retrieve order", http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderRepo.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var order model.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if order.CustomerID == "" {
		http.Error(w, "CustomerID is required", http.StatusBadRequest)
		return
	}
	if len(order.Products) == 0 {
		http.Error(w, "At least one product is required", http.StatusBadRequest)
		return
	}
	if order.TotalAmount <= 0 {
		http.Error(w, "TotalAmount must be greater than zero", http.StatusBadRequest)
		return
	}
	if order.ShippingAddress == "" {
		http.Error(w, "ShippingAddress is required", http.StatusBadRequest)
		return
	}
	if order.PaymentMethod == "" {
		http.Error(w, "PaymentMethod is required", http.StatusBadRequest)
		return
	}

	order.UpdatedAt = time.Now()

	err = h.orderRepo.UpdateOrder(id, &order)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Order updated successfully",
		"order":   order,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	err := h.orderRepo.DeleteOrder(id)
	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Order with id " + id + " deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}
