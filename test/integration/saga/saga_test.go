package saga_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cartClient "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/client"
	orderClient "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/client"
	orderEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	paymentClient "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/client"
	paymentRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
	sagaRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

func TestOrderPaymentSaga(t *testing.T) {
	// Setup test database
	tdb := testutil.NewTestDB(t)
	defer tdb.Cleanup()

	// Clean up tables before test
	require.NoError(t, tdb.TruncateTables("orders", "payments", "saga_transactions", "saga_steps"))

	// Initialize repositories
	sagaRepository := sagaRepo.NewSagaRepository(tdb.DB)
	orderRepository := orderRepo.NewOrderRepository(tdb.DB)
	paymentRepository := paymentRepo.NewPaymentRepository(tdb.DB)

	// Initialize mock gRPC clients
	orderGrpcClient, err := orderClient.NewOrderClient("localhost:50051")
	require.NoError(t, err)
	defer orderGrpcClient.Close()

	paymentGrpcClient, err := paymentClient.NewPaymentClient("localhost:50052")
	require.NoError(t, err)
	defer paymentGrpcClient.Close()

	cartGrpcClient, err := cartClient.NewCartClient("localhost:50053")
	require.NoError(t, err)
	defer cartGrpcClient.Close()

	// Initialize saga usecase
	sagaUsecase := usecase.NewSagaUsecase(
		sagaRepository,
		orderRepository,
		paymentRepository,
		orderGrpcClient,
		paymentGrpcClient,
		cartGrpcClient,
	)

	// Test cases
	t.Run("successful order payment saga", func(t *testing.T) {
		// Create test order
		orderID := uuid.New()
		order := &orderEntity.Order{
			ID:          orderID,
			UserID:      uuid.New(),
			TotalAmount: 100.0,
			Status:      orderEntity.OrderStatusPending,
		}
		err := orderRepository.Create(context.Background(), order)
		require.NoError(t, err)

		// Start saga
		err = sagaUsecase.StartOrderPaymentSaga(context.Background(), orderID)
		require.NoError(t, err)

		// Wait for saga to complete
		time.Sleep(5 * time.Second)

		// Verify saga status
		saga, err := sagaRepository.GetByOrderID(context.Background(), orderID)
		require.NoError(t, err)
		assert.Equal(t, entity.SagaStatusCompleted, saga.Status)

		// Verify all steps completed
		for _, step := range saga.Steps {
			assert.Equal(t, entity.StepStatusCompleted, step.Status)
		}

		// Verify order and payment status
		updatedOrder, err := orderRepository.GetByID(context.Background(), orderID)
		require.NoError(t, err)
		assert.Equal(t, orderEntity.OrderStatusCompleted, updatedOrder.Status)

		payment, err := paymentRepository.GetByOrderID(context.Background(), orderID)
		require.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "SUCCESS", string(payment.Status))
	})

	t.Run("failed payment saga with compensation", func(t *testing.T) {
		// Create test order with invalid amount to trigger failure
		orderID := uuid.New()
		order := &orderEntity.Order{
			ID:          orderID,
			UserID:      uuid.New(),
			TotalAmount: -1.0, // Invalid amount to trigger failure
			Status:      orderEntity.OrderStatusPending,
		}
		err := orderRepository.Create(context.Background(), order)
		require.NoError(t, err)

		// Start saga
		err = sagaUsecase.StartOrderPaymentSaga(context.Background(), orderID)
		require.NoError(t, err)

		// Wait for saga to complete/fail
		time.Sleep(5 * time.Second)

		// Verify saga status
		saga, err := sagaRepository.GetByOrderID(context.Background(), orderID)
		require.NoError(t, err)
		assert.Equal(t, entity.SagaStatusFailed, saga.Status)

		// Verify compensation occurred
		var foundCompensated bool
		for _, step := range saga.Steps {
			if step.Status == entity.StepStatusCompensated {
				foundCompensated = true
				break
			}
		}
		assert.True(t, foundCompensated, "Expected to find at least one compensated step")

		// Verify order status reverted
		updatedOrder, err := orderRepository.GetByID(context.Background(), orderID)
		require.NoError(t, err)
		assert.Equal(t, orderEntity.OrderStatusFailed, updatedOrder.Status)
	})
}
