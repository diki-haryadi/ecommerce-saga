package http

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all payment-related routes
func RegisterRoutes(router fiber.Router, handler *PaymentHandler, authMiddleware fiber.Handler) {
	payments := router.Group("/payments")

	// Protected routes
	protected := payments.Group("")
	protected.Use(authMiddleware)

	protected.Post("", handler.ProcessPayment)
	protected.Get("/:id", handler.GetPayment)
	protected.Get("/order/:id", handler.GetPaymentByOrder)
	protected.Put("/:id/status", handler.UpdatePaymentStatus)

	// Webhook routes (unprotected)
	payments.Post("/webhook/:provider", handler.HandleWebhook)
}
