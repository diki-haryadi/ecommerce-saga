package entity

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents an item in the cart
type CartItem struct {
	CartID    uuid.UUID `json:"cart_id" bson:"cart_id"`
	ProductID uuid.UUID `json:"product_id" bson:"product_id"`
	Name      string    `json:"name" bson:"name"`
	Price     float64   `json:"price" bson:"price"`
	Quantity  int       `json:"quantity" bson:"quantity"`
}

// Cart represents a user's shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id" bson:"_id"`
	UserID    uuid.UUID  `json:"user_id" bson:"user_id"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"total"`
	ExpiresAt time.Time  `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}

// NewCart creates a new cart for a user
func NewCart(userID uuid.UUID, expiry time.Duration) *Cart {
	now := time.Now()
	return &Cart{
		ID:        uuid.New(),
		UserID:    userID,
		Items:     make([]CartItem, 0),
		Total:     0,
		ExpiresAt: now.Add(expiry),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddItem adds a new item to the cart or updates its quantity if it already exists
func (c *Cart) AddItem(item CartItem) {
	for i, existingItem := range c.Items {
		if existingItem.ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			c.calculateTotal()
			return
		}
	}
	c.Items = append(c.Items, item)
	c.calculateTotal()
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID uuid.UUID) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.calculateTotal()
			return
		}
	}
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (c *Cart) UpdateItemQuantity(productID uuid.UUID, quantity int) bool {
	if quantity <= 0 {
		c.RemoveItem(productID)
		return true
	}

	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity = quantity
			c.calculateTotal()
			return true
		}
	}
	return false
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = make([]CartItem, 0)
	c.Total = 0
}

// IsExpired checks if the cart has expired
func (c *Cart) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// calculateTotal recalculates the total price of all items in the cart
func (c *Cart) calculateTotal() {
	var total float64
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	c.Total = total
	c.UpdatedAt = time.Now()
}
