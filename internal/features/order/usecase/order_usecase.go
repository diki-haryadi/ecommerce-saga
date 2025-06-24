package usecase

import (
	"context"
	"errors"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/usecase"

	"github.com/google/uuid"

	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartEmpty         = errors.New("cart is empty")
	ErrInvalidStatus     = errors.New("invalid order status")
	ErrStatusTransition  = errors.New("invalid status transition")
	ErrOrderAlreadyFinal = errors.New("order is already in final state")

	ErrNotFound  = errors.New("order not found")
	ErrCancelled = errors.New("order is already cancelled")
	ErrCompleted = errors.New("order is already completed")
)

type OrderUsecase struct {
	orderRepo repository.OrderRepository
	cartRepo  cartRepo.CartRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, cartRepo cartRepo.CartRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

// CreateOrder creates a new order from a cart
func (u *OrderUsecase) CreateOrder(ctx context.Context, userID, cartID uuid.UUID, paymentMethod, shippingAddress string) (*usecase.OrderResponse, error) {
	// Get cart
	cart, err := u.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, usecase.ErrCartNotFound
	}

	// Validate cart
	if len(cart.Items) == 0 {
		return nil, usecase.ErrCartEmpty
	}
	if cart.UserID != userID {
		return nil, usecase.ErrCartNotFound
	}

	// Create order items
	items := make([]entity.OrderItem, len(cart.Items))
	for i, cartItem := range cart.Items {
		items[i] = entity.OrderItem{
			ID:        uuid.New(),
			ProductID: cartItem.ProductID,
			Name:      cartItem.Name,
			Price:     cartItem.Price,
			Quantity:  cartItem.Quantity,
		}
	}

	// Create order
	newOrder := entity.NewOrder(userID, items)

	// Save order
	if err := u.orderRepo.Create(ctx, newOrder); err != nil {
		return nil, err
	}

	// Clear cart
	if err := u.cartRepo.Delete(ctx, cart.ID); err != nil {
		return nil, err
	}

	// Convert to response
	return &usecase.OrderResponse{
		ID:          newOrder.ID,
		UserID:      newOrder.UserID,
		Items:       u.convertItems(items),
		TotalAmount: newOrder.TotalAmount,
		Status:      usecase.Status(newOrder.Status),
		CreatedAt:   newOrder.CreatedAt,
		UpdatedAt:   newOrder.UpdatedAt,
	}, nil
}

func (u *OrderUsecase) convertItems(items []entity.OrderItem) []usecase.OrderItem {
	result := make([]usecase.OrderItem, len(items))
	for i, item := range items {
		result[i] = usecase.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Price * float64(item.Quantity),
		}
	}
	return result
}

// GetOrder retrieves an order by ID
func (u *OrderUsecase) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*usecase.OrderResponse, error) {
	orderEntity, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if orderEntity == nil || orderEntity.UserID != userID {
		return nil, usecase.ErrNotFound
	}

	return &usecase.OrderResponse{
		ID:          orderEntity.ID,
		UserID:      orderEntity.UserID,
		Items:       u.convertItems(orderEntity.Items),
		TotalAmount: orderEntity.TotalAmount,
		Status:      usecase.Status(orderEntity.Status),
		CreatedAt:   orderEntity.CreatedAt,
		UpdatedAt:   orderEntity.UpdatedAt,
	}, nil
}

// ListOrders retrieves a paginated list of orders for a user
func (u *OrderUsecase) ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*usecase.OrderResponse, int64, error) {
	// Get total count
	totalRows, err := u.orderRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	offset := (page - 1) * limit

	// Get orders
	orders, err := u.orderRepo.GetByUserID(ctx, userID, int(limit), int(offset))
	if err != nil {
		return nil, 0, err
	}

	// Convert to response
	result := make([]*usecase.OrderResponse, len(orders))
	for i, o := range orders {
		result[i] = &usecase.OrderResponse{
			ID:          o.ID,
			UserID:      o.UserID,
			Items:       u.convertItems(o.Items),
			TotalAmount: o.TotalAmount,
			Status:      usecase.Status(o.Status),
			CreatedAt:   o.CreatedAt,
			UpdatedAt:   o.UpdatedAt,
		}
	}

	return result, totalRows, nil
}

// UpdateOrderStatus updates the status of an order
func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status usecase.Status) (*usecase.OrderResponse, error) {
	// Get order
	orderEntity, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if orderEntity == nil {
		return nil, usecase.ErrNotFound
	}

	// Validate status transition
	newStatus := entity.OrderStatus(status)
	if !orderEntity.CanTransitionTo(newStatus) {
		if orderEntity.IsFinal() {
			return nil, usecase.ErrOrderAlreadyFinal
		}
		return nil, usecase.ErrStatusTransition
	}

	// Update status
	if err := u.orderRepo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		return nil, err
	}

	// Get updated order
	orderEntity, err = u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &usecase.OrderResponse{
		ID:          orderEntity.ID,
		UserID:      orderEntity.UserID,
		Items:       u.convertItems(orderEntity.Items),
		TotalAmount: orderEntity.TotalAmount,
		Status:      usecase.Status(orderEntity.Status),
		CreatedAt:   orderEntity.CreatedAt,
		UpdatedAt:   orderEntity.UpdatedAt,
	}, nil
}

// CancelOrder cancels an order if possible
func (u *OrderUsecase) CancelOrder(ctx context.Context, userID, orderID uuid.UUID, reason string) error {
	// Get order
	orderEntity, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if orderEntity == nil || orderEntity.UserID != userID {
		return usecase.ErrNotFound
	}

	// Check if order can be cancelled
	if orderEntity.Status == entity.OrderStatus(usecase.StatusCancelled) {
		return usecase.ErrCancelled
	}
	if orderEntity.Status == entity.OrderStatus(usecase.StatusDelivered) {
		return usecase.ErrCompleted
	}

	// Update status to cancelled
	if err := u.orderRepo.UpdateStatus(ctx, orderID, entity.OrderStatus(usecase.StatusCancelled)); err != nil {
		return err
	}

	return nil
}
