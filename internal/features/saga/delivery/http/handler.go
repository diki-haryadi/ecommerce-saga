package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/errors"
	httpresponse "github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
)

type SagaHandler struct {
	sagaUsecase  *usecase.SagaUsecase
	errorHandler errors.ErrorHandler
}

func NewSagaHandler(sagaUsecase *usecase.SagaUsecase) *SagaHandler {
	return &SagaHandler{
		sagaUsecase:  sagaUsecase,
		errorHandler: errors.NewErrorHandler(),
	}
}

// StartOrderPaymentSaga handles POST /sagas/order-payment request
func (h *SagaHandler) StartOrderPaymentSaga(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Query("order_id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid order ID"))
	}

	if err := h.sagaUsecase.StartOrderPaymentSaga(c.Context(), orderID); err != nil {
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Saga started successfully", nil)
}

// GetSagaStatus handles GET /sagas/:id request
func (h *SagaHandler) GetSagaStatus(c *fiber.Ctx) error {
	sagaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid saga ID"))
	}

	status, err := h.sagaUsecase.GetSagaStatus(c.Context(), sagaID)
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Saga status retrieved successfully", status)
}
