package request

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	CartID          string `json:"cart_id" validate:"required"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	ShippingAddress string `json:"shipping_address" validate:"required"`
}

// UpdateOrderStatusRequest represents the request to update an order's status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID SHIPPED DELIVERED CANCELLED"`
}

// ListOrdersRequest represents the request to list user orders
type ListOrdersRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Status   string `form:"status" validate:"omitempty,oneof=PENDING PAID SHIPPED DELIVERED CANCELLED"`
}
