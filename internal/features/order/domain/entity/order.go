package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusFailed     OrderStatus = "FAILED"
)

// OrderItem represents an item in the order
type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
}

// Order represents an order in the system
type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewOrder creates a new order from cart items
func NewOrder(userID uuid.UUID, items []OrderItem) *Order {
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	return &Order{
		ID:          uuid.New(),
		UserID:      userID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// IsFinal checks if the order is in a final state
func (o *Order) IsFinal() bool {
	return o.Status == OrderStatusDelivered ||
		o.Status == OrderStatusCancelled ||
		o.Status == OrderStatusFailed
}

// CanTransitionTo checks if the order can transition to the given status
func (o *Order) CanTransitionTo(status OrderStatus) bool {
	if o.IsFinal() {
		return false
	}

	switch o.Status {
	case OrderStatusPending:
		return status == OrderStatusConfirmed || status == OrderStatusCancelled
	case OrderStatusConfirmed:
		return status == OrderStatusProcessing || status == OrderStatusCancelled
	case OrderStatusProcessing:
		return status == OrderStatusShipped || status == OrderStatusFailed
	case OrderStatusShipped:
		return status == OrderStatusDelivered || status == OrderStatusFailed
	default:
		return false
	}
}
