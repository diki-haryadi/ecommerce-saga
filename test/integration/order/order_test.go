package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

func TestOrderRepository(t *testing.T) {
	// Setup test database
	tdb := testutil.NewTestDB(t)
	defer tdb.Cleanup()

	// Clean up tables before test
	require.NoError(t, tdb.TruncateTables("orders"))

	// Initialize repository
	orderRepo := postgres.NewOrderRepository(tdb.DB)

	t.Run("create and retrieve order", func(t *testing.T) {
		// Create test order
		order := &entity.Order{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Status: entity.OrderStatusPending,
			Items: []entity.OrderItem{
				{
					ProductID: uuid.New(),
					Quantity:  2,
					Price:     50.0,
				},
			},
		}

		// Create order
		err := orderRepo.Create(context.Background(), order)
		require.NoError(t, err)

		// Retrieve order
		retrieved, err := orderRepo.GetByID(context.Background(), order.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, order.ID, retrieved.ID)
		assert.Equal(t, order.UserID, retrieved.UserID)
		assert.Equal(t, order.Status, retrieved.Status)
		assert.Len(t, retrieved.Items, 1)
	})

	t.Run("update order status", func(t *testing.T) {
		// Create test order
		order := &entity.Order{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Status: entity.OrderStatusPending,
		}

		// Create order
		err := orderRepo.Create(context.Background(), order)
		require.NoError(t, err)

		// Update status
		order.Status = entity.OrderStatusCompleted
		err = orderRepo.Update(context.Background(), order)
		require.NoError(t, err)

		// Verify status update
		updated, err := orderRepo.GetByID(context.Background(), order.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.OrderStatusCompleted, updated.Status)
	})

	t.Run("list orders by user", func(t *testing.T) {
		userID := uuid.New()

		// Create multiple orders for the same user
		for i := 0; i < 3; i++ {
			order := &entity.Order{
				ID:     uuid.New(),
				UserID: userID,
				Status: entity.OrderStatusPending,
			}
			err := orderRepo.Create(context.Background(), order)
			require.NoError(t, err)
		}

		// List orders
		orders, err := orderRepo.ListByUserID(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, orders, 3)

		// Verify all orders belong to the user
		for _, order := range orders {
			assert.Equal(t, userID, order.UserID)
		}
	})
}
