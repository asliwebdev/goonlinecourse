package handler

import (
	"net/http"

	"lesson22.5/repository"

	"github.com/gorilla/mux"
)

type Handler struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	orderRepo    *repository.OrderRepository
}

func NewHandler(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository, orderRepo *repository.OrderRepository) *Handler {
	return &Handler{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		orderRepo:    orderRepo,
	}
}

func Run(handler *Handler) *http.Server {
	router := mux.NewRouter()

	// PRODUCT ROUTES
	router.HandleFunc("/product", handler.CreateProduct).Methods("POST")
	router.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", handler.GetProductById).Methods("GET")
	router.HandleFunc("/product/{id}", handler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/{id}", handler.DeleteProduct).Methods("DELETE")

	// CATEGORY ROUTES
	router.HandleFunc("/category", handler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories", handler.GetAllCategories).Methods("GET")
	router.HandleFunc("/category/{id}", handler.GetCategoryById).Methods("GET")
	router.HandleFunc("/category/{id}", handler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/category/{id}", handler.DeleteCategory).Methods("DELETE")

	// ORDER ROUTES
	router.HandleFunc("/order", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders", handler.GetAllOrders).Methods("GET")
	router.HandleFunc("/order/{id}", handler.GetOrderById).Methods("GET")
	router.HandleFunc("/order/{id}", handler.UpdateOrder).Methods("PUT")
	router.HandleFunc("/order/{id}", handler.DeleteOrder).Methods("DELETE")

	server := &http.Server{Addr: ":8080", Handler: router}

	return server
}
