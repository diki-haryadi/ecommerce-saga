package http

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all cart-related routes
func RegisterRoutes(router fiber.Router, handler *CartHandler, authMiddleware fiber.Handler) {
	cart := router.Group("/cart")
	cart.Use(authMiddleware)

	cart.Get("", handler.GetCart)
	cart.Post("/items", handler.AddItem)
	cart.Put("/items/:id", handler.UpdateItem)
	cart.Delete("/items/:id", handler.RemoveItem)
	cart.Delete("", handler.ClearCart)
}
