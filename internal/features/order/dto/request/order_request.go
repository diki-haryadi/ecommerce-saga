package request

import "github.com/google/uuid"

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	CartID uuid.UUID `json:"cart_id" validate:"required"`
}

// UpdateOrderStatusRequest represents the request to update an order's status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=CONFIRMED PROCESSING SHIPPED DELIVERED CANCELLED FAILED"`
}

// ListOrdersRequest represents the request to list user orders
type ListOrdersRequest struct {
	Page     int `form:"page" validate:"min=1"`
	PageSize int `form:"page_size" validate:"min=1,max=100"`
}
