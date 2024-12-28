package handler

import (
	"lesson23/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderRepo *repository.OrderRepository
}

func NewHandler(orderRepo *repository.OrderRepository) *Handler {
	return &Handler{
		orderRepo: orderRepo,
	}
}

func Run(h *Handler) *gin.Engine {
	router := gin.Default()

	// ORDER ROUTES
	orderRoutes := router.Group("/orders")
	{
		orderRoutes.POST("/", h.CreateOrder)
		orderRoutes.GET("/", h.GetAllOrders)
		orderRoutes.GET("/:id", h.GetOrderById)
		orderRoutes.PUT("/:id", h.UpdateOrder)
		orderRoutes.DELETE("/:id", h.DeleteOrder)
	}

	return router
}
