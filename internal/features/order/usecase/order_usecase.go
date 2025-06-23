package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/dto/response"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartEmpty         = errors.New("cart is empty")
	ErrInvalidStatus     = errors.New("invalid order status")
	ErrStatusTransition  = errors.New("invalid status transition")
	ErrOrderAlreadyFinal = errors.New("order is already in final state")
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
func (u *OrderUsecase) CreateOrder(ctx context.Context, userID uuid.UUID, req *request.CreateOrderRequest) (*response.OrderResponse, error) {
	// Get cart
	cart, err := u.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Validate cart
	if len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}
	if cart.UserID != userID {
		return nil, ErrCartNotFound
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
	order := entity.NewOrder(userID, items)

	// Save order
	if err := u.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Clear cart
	if err := u.cartRepo.Delete(ctx, cart.ID); err != nil {
		return nil, err
	}

	return response.NewOrderResponse(order), nil
}

// GetOrder retrieves an order by ID
func (u *OrderUsecase) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*response.OrderResponse, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil || order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	return response.NewOrderResponse(order), nil
}

// ListOrders retrieves a paginated list of orders for a user
func (u *OrderUsecase) ListOrders(ctx context.Context, userID uuid.UUID, req *request.ListOrdersRequest) (*response.OrderListResponse, error) {
	// Get total count
	totalRows, err := u.orderRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (req.Page - 1) * req.PageSize

	// Get orders
	orders, err := u.orderRepo.GetByUserID(ctx, userID, req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	return response.NewOrderListResponse(orders, req.Page, req.PageSize, totalRows), nil
}

// UpdateOrderStatus updates the status of an order
func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, userID, orderID uuid.UUID, req *request.UpdateOrderStatusRequest) (*response.OrderResponse, error) {
	// Get order
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil || order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	// Validate status transition
	newStatus := entity.OrderStatus(req.Status)
	if !order.CanTransitionTo(newStatus) {
		if order.IsFinal() {
			return nil, ErrOrderAlreadyFinal
		}
		return nil, ErrStatusTransition
	}

	// Update status
	if err := u.orderRepo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		return nil, err
	}

	// Get updated order
	order, err = u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return response.NewOrderResponse(order), nil
}
