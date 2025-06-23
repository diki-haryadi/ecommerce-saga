package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/usecase"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

// ProcessPayment handles POST /payments request
func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	var req request.ProcessPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	payment, err := h.paymentUsecase.ProcessPayment(c.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrInvalidProvider:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process payment"})
		}
	}

	return c.Status(http.StatusAccepted).JSON(payment)
}

// GetPayment handles GET /payments/:id request
func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	payment, err := h.paymentUsecase.GetPayment(c.Context(), id)
	if err != nil {
		if err == usecase.ErrPaymentNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get payment"})
	}

	return c.Status(http.StatusOK).JSON(payment)
}

// GetPaymentByOrder handles GET /payments/order/:id request
func (h *PaymentHandler) GetPaymentByOrder(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	payment, err := h.paymentUsecase.GetPaymentByOrder(c.Context(), orderID)
	if err != nil {
		if err == usecase.ErrPaymentNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get payment"})
	}

	return c.Status(http.StatusOK).JSON(payment)
}

// UpdatePaymentStatus handles PUT /payments/:id/status request
func (h *PaymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	var req request.UpdatePaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	payment, err := h.paymentUsecase.UpdatePaymentStatus(c.Context(), id, &req)
	if err != nil {
		switch err {
		case usecase.ErrPaymentNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrStatusTransition, usecase.ErrPaymentCompleted:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update payment status"})
		}
	}

	return c.Status(http.StatusOK).JSON(payment)
}

// HandleWebhook handles POST /payments/webhook/:provider request
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	provider := entity.PaymentProvider(c.Params("provider"))
	signature := c.Get("X-Webhook-Signature")

	payload := c.Body()

	if err := h.paymentUsecase.HandleWebhook(c.Context(), provider, payload, signature); err != nil {
		switch err {
		case usecase.ErrInvalidProvider:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process webhook"})
		}
	}

	return c.SendStatus(http.StatusOK)
}
