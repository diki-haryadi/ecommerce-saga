package http

import (
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers payment routes
func RegisterRoutes(router fiber.Router, useCase usecase.Usecase) {
	handler := NewPaymentHandler(useCase)

	paymentGroup := router.Group("/payments")
	{
		paymentGroup.Post("/", handler.CreatePayment)
		paymentGroup.Get("/:id", handler.GetPayment)
		paymentGroup.Get("/", handler.ListPayments)
		paymentGroup.Post("/:id/process", handler.ProcessPayment)
		paymentGroup.Post("/:id/refund", handler.RefundPayment)
	}
}
