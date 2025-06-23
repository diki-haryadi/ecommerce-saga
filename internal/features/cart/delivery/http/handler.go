package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/errors"
	httpresponse "github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
)

type CartHandler struct {
	cartUsecase  *usecase.CartUsecase
	errorHandler errors.ErrorHandler
}

func NewCartHandler(cartUsecase *usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase:  cartUsecase,
		errorHandler: errors.NewErrorHandler(),
	}
}

// GetCart handles GET /cart request
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	resp, err := h.cartUsecase.GetCart(c.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrCartExpired:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Cart retrieved successfully", resp)
}

// AddItem handles POST /cart/items request
func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	var req request.AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	resp, err := h.cartUsecase.AddItem(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartExpired:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		case usecase.ErrProductNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Item added to cart successfully", resp)
}

// UpdateItem handles PUT /cart/items/:id request
func (h *CartHandler) UpdateItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	var req request.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	resp, err := h.cartUsecase.UpdateItem(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrCartExpired:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		case usecase.ErrItemNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Cart item updated successfully", resp)
}

// RemoveItem handles DELETE /cart/items/:id request
func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid product ID"))
	}

	req := &request.RemoveItemRequest{
		ProductID: productID,
	}

	resp, err := h.cartUsecase.RemoveItem(c.Context(), userID, req)
	if err != nil {
		switch err {
		case usecase.ErrCartNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrCartExpired:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Item removed from cart successfully", resp)
}

// ClearCart handles DELETE /cart request
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid user ID"))
	}

	if err := h.cartUsecase.ClearCart(c.Context(), userID); err != nil {
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "Cart cleared successfully", nil)
}
