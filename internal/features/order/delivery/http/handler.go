package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/errors"
	httpresponse "github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
	errorHandler errors.ErrorHandler
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
		errorHandler: errors.NewErrorHandler(),
	}
}

// CreateOrder handles POST /orders request
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	var req request.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	resp, err := h.orderUsecase.CreateOrder(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrCartEmpty:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.Created(c, "Order created successfully", resp)
}

// GetOrder handles GET /orders/:id request
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid order ID"))
	}

	resp, err := h.orderUsecase.GetOrder(c.Context(), userID, orderID)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Order retrieved successfully", resp)
}

// ListOrders handles GET /orders request
func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	var req request.ListOrdersRequest
	if err := c.QueryParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	// Set default values if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resp, err := h.orderUsecase.ListOrders(c.Context(), userID, &req)
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Orders retrieved successfully", resp)
}

// UpdateOrderStatus handles PUT /orders/:id/status request
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid order ID"))
	}

	var req request.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	resp, err := h.orderUsecase.UpdateOrderStatus(c.Context(), userID, orderID, &req)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrStatusTransition:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		case usecase.ErrOrderAlreadyFinal:
			return h.errorHandler.Handle(c, errors.NewConflictError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Order status updated successfully", resp)
}
