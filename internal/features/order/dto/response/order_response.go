package response

import (
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
)

// OrderItemResponse represents an order item in responses
type OrderItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
}

// OrderResponse represents an order in responses
type OrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	Orders    []OrderResponse `json:"orders"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalRows int64           `json:"total_rows"`
}

// NewOrderResponse creates a new order response from an order entity
func NewOrderResponse(order *entity.Order) *OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Price * float64(item.Quantity),
		}
	}

	return &OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Items:       items,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

// NewOrderListResponse creates a new order list response
func NewOrderListResponse(orders []*entity.Order, page, pageSize int, totalRows int64) *OrderListResponse {
	orderResponses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		resp := NewOrderResponse(order)
		orderResponses[i] = *resp
	}

	return &OrderListResponse{
		Orders:    orderResponses,
		Page:      page,
		PageSize:  pageSize,
		TotalRows: totalRows,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
