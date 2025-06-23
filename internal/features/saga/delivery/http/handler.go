package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
)

type SagaHandler struct {
	sagaUsecase *usecase.SagaUsecase
}

func NewSagaHandler(sagaUsecase *usecase.SagaUsecase) *SagaHandler {
	return &SagaHandler{
		sagaUsecase: sagaUsecase,
	}
}

// StartOrderPaymentSaga handles POST /sagas/order-payment request
func (h *SagaHandler) StartOrderPaymentSaga(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Query("order_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	if err := h.sagaUsecase.StartOrderPaymentSaga(c.Context(), orderID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start saga"})
	}

	return c.SendStatus(http.StatusAccepted)
}
