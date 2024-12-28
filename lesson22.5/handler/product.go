package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"lesson22.5/model"
)

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var product model.Product

	err = json.Unmarshal(body, &product)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "Product price must be greater than 0", http.StatusBadRequest)
		return
	}

	if product.Stock < 0 {
		http.Error(w, "Product stock cannot be negative", http.StatusBadRequest)
		return
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	err = h.productRepo.CreateProduct(&product)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Product '" + product.Name + "' created successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	product, err := h.productRepo.GetProductById(id)
	if err != nil {
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productRepo.GetAllProducts()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var product model.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "Product price must be greater than 0", http.StatusBadRequest)
		return
	}

	if product.Stock < 0 {
		http.Error(w, "Product stock cannot be negative", http.StatusBadRequest)
		return
	}

	product.UpdatedAt = time.Now()

	err = h.productRepo.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Product updated successfully",
		"product": product,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	err := h.productRepo.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Product with id " + id + " deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}
