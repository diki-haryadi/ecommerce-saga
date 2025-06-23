package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/usecase"
)

type CartHandler struct {
	cartUsecase *usecase.CartUsecase
}

func NewCartHandler(cartUsecase *usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cartUsecase,
	}
}

// GetCart handles GET /cart request
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	resp, err := h.cartUsecase.GetCart(c.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrCartExpired:
			return c.Status(http.StatusGone).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// AddItem handles POST /cart/items request
func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req request.AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	resp, err := h.cartUsecase.AddItem(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartExpired:
			return c.Status(http.StatusGone).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrProductNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// UpdateItem handles PUT /cart/items/:id request
func (h *CartHandler) UpdateItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req request.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	resp, err := h.cartUsecase.UpdateItem(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrCartExpired:
			return c.Status(http.StatusGone).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrItemNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// RemoveItem handles DELETE /cart/items/:id request
func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	req := &request.RemoveItemRequest{
		ProductID: productID,
	}

	resp, err := h.cartUsecase.RemoveItem(c.Context(), userID, req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrCartExpired:
			return c.Status(http.StatusGone).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// ClearCart handles DELETE /cart request
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.cartUsecase.ClearCart(c.Context(), userID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Cart cleared successfully"})
}
