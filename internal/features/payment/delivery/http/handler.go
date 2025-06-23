package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/errors"
	httpresponse "github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
	errorHandler   errors.ErrorHandler
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
		errorHandler:   errors.NewErrorHandler(),
	}
}

// ProcessPayment handles POST /payments request
func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	var req request.ProcessPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	payment, err := h.paymentUsecase.ProcessPayment(c.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrInvalidProvider:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.Created(c, "Payment processed successfully", payment)
}

// GetPayment handles GET /payments/:id request
func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid payment ID"))
	}

	payment, err := h.paymentUsecase.GetPayment(c.Context(), id)
	if err != nil {
		if err == usecase.ErrPaymentNotFound {
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		}
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Payment retrieved successfully", payment)
}

// GetPaymentByOrder handles GET /payments/order/:id request
func (h *PaymentHandler) GetPaymentByOrder(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid order ID"))
	}

	payment, err := h.paymentUsecase.GetPaymentByOrder(c.Context(), orderID)
	if err != nil {
		if err == usecase.ErrPaymentNotFound {
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		}
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Payment retrieved successfully", payment)
}

// UpdatePaymentStatus handles PUT /payments/:id/status request
func (h *PaymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid payment ID"))
	}

	var req request.UpdatePaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	payment, err := h.paymentUsecase.UpdatePaymentStatus(c.Context(), id, &req)
	if err != nil {
		switch err {
		case usecase.ErrPaymentNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrStatusTransition, usecase.ErrPaymentCompleted:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Payment status updated successfully", payment)
}

// HandleWebhook handles POST /payments/webhook/:provider request
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	provider := entity.PaymentProvider(c.Params("provider"))
	signature := c.Get("X-Webhook-Signature")

	payload := c.Body()

	if err := h.paymentUsecase.HandleWebhook(c.Context(), provider, payload, signature); err != nil {
		switch err {
		case usecase.ErrInvalidProvider:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Webhook processed successfully", nil)
}
