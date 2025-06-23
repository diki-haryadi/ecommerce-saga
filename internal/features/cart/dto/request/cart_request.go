package request

import "github.com/google/uuid"

// AddItemRequest represents the request to add an item to the cart
type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// UpdateItemRequest represents the request to update an item's quantity in the cart
type UpdateItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=0"`
}

// RemoveItemRequest represents the request to remove an item from the cart
type RemoveItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
