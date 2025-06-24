package http

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all saga-related routes
func RegisterRoutes(router fiber.Router, handler *SagaHandler, authMiddleware fiber.Handler) {
	sagas := router.Group("/sagas")
	sagas.Use(authMiddleware)

	sagas.Post("/order-payment", handler.StartOrderPaymentSaga)
	sagas.Get("/:id", handler.GetSagaStatus)
}
