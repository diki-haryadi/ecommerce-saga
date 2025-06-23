package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cartdomain "github.com/diki-haryadi/ecommerce-saga/internal/features/cart"
	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/proto"
)

type CartServer struct {
	pb.UnimplementedCartServiceServer
	cartUsecase cartdomain.Usecase
}

func NewCartServer(cartUsecase cartdomain.Usecase) *CartServer {
	return &CartServer{
		cartUsecase: cartUsecase,
	}
}

// toProtoCart converts a domain cart to a protobuf cart
func toProtoCart(cart *cartdomain.Cart) *pb.Cart {
	if cart == nil {
		return nil
	}

	items := make([]*pb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = &pb.CartItem{
			Id:        item.ID.String(),
			ProductId: item.ProductID.String(),
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &pb.Cart{
		Id:     cart.ID.String(),
		UserId: cart.UserID.String(),
		Items:  items,
		Total:  cart.Total,
	}
}

func (s *CartServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	cart, err := s.cartUsecase.AddItem(ctx, userID, productID, req.Quantity)
	if err != nil {
		switch err {
		case cartdomain.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case cartdomain.ErrInvalidQuantity:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case cartdomain.ErrOutOfStock:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to add item to cart")
		}
	}

	return &pb.AddItemResponse{
		Success: true,
		Message: "Item added to cart successfully",
		Cart:    toProtoCart(cart),
	}, nil
}

func (s *CartServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.RemoveItemResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	cartItemID, err := uuid.Parse(req.CartItemId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid cart item ID")
	}

	err = s.cartUsecase.RemoveItem(ctx, userID, cartItemID)
	if err != nil {
		switch err {
		case cartdomain.ErrCartNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case cartdomain.ErrCartItemNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to remove item from cart")
		}
	}

	return &pb.RemoveItemResponse{
		Success: true,
		Message: "Item removed from cart successfully",
	}, nil
}

func (s *CartServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	cartItemID, err := uuid.Parse(req.CartItemId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid cart item ID")
	}

	cart, err := s.cartUsecase.UpdateItem(ctx, userID, cartItemID, req.Quantity)
	if err != nil {
		switch err {
		case cartdomain.ErrCartNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case cartdomain.ErrCartItemNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case cartdomain.ErrInvalidQuantity:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case cartdomain.ErrOutOfStock:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update cart item")
		}
	}

	return &pb.UpdateItemResponse{
		Success: true,
		Message: "Cart item updated successfully",
		Cart:    toProtoCart(cart),
	}, nil
}

func (s *CartServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	cart, err := s.cartUsecase.GetCart(ctx, userID)
	if err != nil {
		switch err {
		case cartdomain.ErrCartNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get cart")
		}
	}

	return &pb.GetCartResponse{
		Cart: toProtoCart(cart),
	}, nil
}

func (s *CartServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.ClearCartResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	err = s.cartUsecase.ClearCart(ctx, userID)
	if err != nil {
		switch err {
		case cartdomain.ErrCartNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to clear cart")
		}
	}

	return &pb.ClearCartResponse{
		Success: true,
		Message: "Cart cleared successfully",
	}, nil
}
