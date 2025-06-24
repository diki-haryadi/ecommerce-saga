package http

import (
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	useCase usecase.Usecase
}

func NewPaymentHandler(useCase usecase.Usecase) *PaymentHandler {
	return &PaymentHandler{
		useCase: useCase,
	}
}

type CreatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Currency      string    `json:"currency" validate:"required,len=3"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
}

type ProcessPaymentRequest struct {
	CardNumber  string `json:"card_number" validate:"required,creditcard"`
	ExpiryMonth string `json:"expiry_month" validate:"required,len=2"`
	ExpiryYear  string `json:"expiry_year" validate:"required,len=2"`
	CVV         string `json:"cvv" validate:"required,len=3"`
	HolderName  string `json:"holder_name" validate:"required"`
}

type RefundPaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Reason string  `json:"reason" validate:"required"`
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payment, err := h.useCase.CreatePayment(c.Context(), req.OrderID, req.Amount, req.Currency, req.PaymentMethod)
	if err != nil {
		if err == payment.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	payment, err := h.useCase.GetPayment(c.Context(), id)
	if err != nil {
		if err == payment.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get payment",
		})
	}

	return c.JSON(payment)
}

func (h *PaymentHandler) ListPayments(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status")

	payments, total, err := h.useCase.ListPayments(c.Context(), userID, int32(page), int32(limit), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list payments",
		})
	}

	return c.JSON(fiber.Map{
		"data": payments,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	var req ProcessPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	details := &usecase.PaymentDetails{
		CardNumber:  req.CardNumber,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
		HolderName:  req.HolderName,
	}

	payment, err := h.useCase.ProcessPayment(c.Context(), id, details)
	if err != nil {
		if err == payment.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == payment.ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process payment",
		})
	}

	return c.JSON(payment)
}

func (h *PaymentHandler) RefundPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment ID",
		})
	}

	var req RefundPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payment, reason, err := h.useCase.RefundPayment(c.Context(), id, req.Amount, req.Reason)
	if err != nil {
		if err == payment.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == payment.ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refund payment",
		})
	}

	return c.JSON(fiber.Map{
		"payment": payment,
		"reason":  reason,
	})
}
