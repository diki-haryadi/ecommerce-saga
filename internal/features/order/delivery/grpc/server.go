package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order"
	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/proto"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	orderUsecase order.Usecase
}

func NewOrderServer(orderUsecase order.Usecase) *OrderServer {
	return &OrderServer{
		orderUsecase: orderUsecase,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	cartID, err := uuid.Parse(req.CartId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid cart ID")
	}

	orderResp, err := s.orderUsecase.CreateOrder(ctx, userID, cartID, req.PaymentMethod, req.ShippingAddress)
	if err != nil {
		switch err {
		case order.ErrCartNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case order.ErrCartEmpty:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to create order")
		}
	}

	return &pb.CreateOrderResponse{
		Success: true,
		Message: "Order created successfully",
		Order:   convertOrderToPb(orderResp),
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	orderResp, err := s.orderUsecase.GetOrder(ctx, userID, orderID)
	if err != nil {
		switch err {
		case order.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get order")
		}
	}

	return &pb.GetOrderResponse{
		Order: convertOrderToPb(orderResp),
	}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	ordersResp, total, err := s.orderUsecase.ListOrders(ctx, userID, req.Page, req.Limit, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list orders")
	}

	orders := make([]*pb.Order, len(ordersResp))
	for i, o := range ordersResp {
		orders[i] = convertOrderToPb(o)
	}

	return &pb.ListOrdersResponse{
		Orders: orders,
		Total:  total,
		Page:   req.Page,
		Limit:  req.Limit,
	}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	err = s.orderUsecase.CancelOrder(ctx, userID, orderID, req.Reason)
	if err != nil {
		switch err {
		case order.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case order.ErrCancelled:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case order.ErrCompleted:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case order.ErrOrderAlreadyFinal:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to cancel order")
		}
	}

	return &pb.CancelOrderResponse{
		Success: true,
		Message: "Order cancelled successfully",
	}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	orderStatus := order.Status(req.Status)
	orderResp, err := s.orderUsecase.UpdateOrderStatus(ctx, orderID, orderStatus)
	if err != nil {
		switch err {
		case order.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case order.ErrInvalidStatus:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case order.ErrStatusTransition:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case order.ErrOrderAlreadyFinal:
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update order status")
		}
	}

	return &pb.UpdateOrderStatusResponse{
		Success: true,
		Message: "Order status updated successfully",
		Order:   convertOrderToPb(orderResp),
	}, nil
}

func convertOrderToPb(order *order.OrderResponse) *pb.Order {
	pbItems := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		pbItems[i] = &pb.OrderItem{
			Id:        item.ID.String(),
			ProductId: item.ProductID.String(),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
			Subtotal:  item.Subtotal,
		}
	}

	return &pb.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Items:       pbItems,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   timestamppb.New(order.CreatedAt),
		UpdatedAt:   timestamppb.New(order.UpdatedAt),
	}
}
