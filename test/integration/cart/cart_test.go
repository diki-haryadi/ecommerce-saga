package cart_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

func TestCartRepository(t *testing.T) {
	// Setup test database
	tdb := testutil.NewTestDB(t)
	defer tdb.Cleanup()

	// Clean up tables before test
	require.NoError(t, tdb.TruncateTables("carts", "cart_items"))

	// Initialize repository
	cartRepo := postgres.NewCartRepository(tdb.DB)

	t.Run("create and retrieve cart", func(t *testing.T) {
		// Create test cart
		userID := uuid.New()
		cart := &entity.Cart{
			ID:     uuid.New(),
			UserID: userID,
			Items: []entity.CartItem{
				{
					ProductID: uuid.New(),
					Quantity:  2,
					Price:     50.0,
				},
			},
		}

		// Create cart
		err := cartRepo.Create(context.Background(), cart)
		require.NoError(t, err)

		// Retrieve cart
		retrieved, err := cartRepo.GetByUserID(context.Background(), userID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, cart.ID, retrieved.ID)
		assert.Equal(t, userID, retrieved.UserID)
		assert.Len(t, retrieved.Items, 1)
	})

	t.Run("add item to cart", func(t *testing.T) {
		// Create test cart
		userID := uuid.New()
		cart := &entity.Cart{
			ID:     uuid.New(),
			UserID: userID,
		}

		// Create cart
		err := cartRepo.Create(context.Background(), cart)
		require.NoError(t, err)

		// Add item
		item := entity.CartItem{
			CartID:    cart.ID,
			ProductID: uuid.New(),
			Quantity:  1,
			Price:     25.0,
		}
		err = cartRepo.AddItem(context.Background(), &item)
		require.NoError(t, err)

		// Verify item added
		updated, err := cartRepo.GetByUserID(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, updated.Items, 1)
		assert.Equal(t, item.ProductID, updated.Items[0].ProductID)
	})

	t.Run("clear cart", func(t *testing.T) {
		// Create test cart with items
		userID := uuid.New()
		cart := &entity.Cart{
			ID:     uuid.New(),
			UserID: userID,
			Items: []entity.CartItem{
				{
					ProductID: uuid.New(),
					Quantity:  2,
					Price:     50.0,
				},
				{
					ProductID: uuid.New(),
					Quantity:  1,
					Price:     30.0,
				},
			},
		}

		// Create cart
		err := cartRepo.Create(context.Background(), cart)
		require.NoError(t, err)

		// Clear cart
		err = cartRepo.Clear(context.Background(), userID)
		require.NoError(t, err)

		// Verify cart is empty
		updated, err := cartRepo.GetByUserID(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, updated.Items, 0)
	})
}
