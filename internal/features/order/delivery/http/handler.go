package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/usecase"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

// CreateOrder handles POST /orders request
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req request.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	resp, err := h.orderUsecase.CreateOrder(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrCartEmpty:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusCreated).JSON(resp)
}

// GetOrder handles GET /orders/:id request
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	resp, err := h.orderUsecase.GetOrder(c.Context(), userID, orderID)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// ListOrders handles GET /orders request
func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req request.ListOrdersRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// UpdateOrderStatus handles PUT /orders/:id/status request
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req request.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	resp, err := h.orderUsecase.UpdateOrderStatus(c.Context(), userID, orderID, &req)
	if err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrStatusTransition:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrOrderAlreadyFinal:
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}
