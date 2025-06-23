package http

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all order-related routes
func RegisterRoutes(router fiber.Router, handler *OrderHandler, authMiddleware fiber.Handler) {
	orders := router.Group("/orders")
	orders.Use(authMiddleware)

	orders.Post("", handler.CreateOrder)
	orders.Get("", handler.ListOrders)
	orders.Get("/:id", handler.GetOrder)
	orders.Put("/:id/status", handler.UpdateOrderStatus)
}
