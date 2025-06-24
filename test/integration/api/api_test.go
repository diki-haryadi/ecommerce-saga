package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/http"
	cartEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository/postgres"
	orderEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	sagaHandler "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/http"
	sagaRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

func setupTestApp(t *testing.T) (*fiber.App, *testutil.TestDB) {
	// Setup database
	tdb := testutil.NewTestDB(t)

	// Initialize repositories
	cartRepository := cartRepo.NewCartRepository(tdb.DB)
	orderRepository := orderRepo.NewOrderRepository(tdb.DB)
	sagaRepository := sagaRepo.NewSagaRepository(tdb.DB)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	api := app.Group("/api")

	cartHandler := http.NewCartHandler(cartRepository)
	cartGroup := api.Group("/cart")
	cartGroup.Post("/", cartHandler.CreateCart)
	cartGroup.Get("/:user_id", cartHandler.GetCart)
	cartGroup.Post("/items", cartHandler.AddItem)
	cartGroup.Delete("/:user_id", cartHandler.ClearCart)

	// Initialize saga usecase and handler
	sagaUsecase := usecase.NewSagaUsecase(sagaRepository, orderRepository, nil, nil, nil, nil)
	sagaHandler := sagaHandler.NewSagaHandler(sagaUsecase)
	sagaGroup := api.Group("/saga")
	sagaGroup.Post("/order-payment", sagaHandler.StartOrderPaymentSaga)
	sagaGroup.Get("/:id", sagaHandler.GetSagaStatus)

	return app, tdb
}

func TestCartAPI(t *testing.T) {
	app, tdb := setupTestApp(t)
	defer tdb.Cleanup()

	t.Run("create and get cart", func(t *testing.T) {
		userID := uuid.New()

		// Create cart
		createReq := fiber.Map{
			"user_id": userID.String(),
		}
		req, _ := json.Marshal(createReq)
		resp, err := app.Test(httptest.NewRequest("POST", "/api/cart", bytes.NewReader(req)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Get cart
		resp, err = app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/cart/%s", userID), nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cart cartEntity.Cart
		err = json.NewDecoder(resp.Body).Decode(&cart)
		require.NoError(t, err)
		assert.Equal(t, userID, cart.UserID)
	})

	t.Run("add item to cart", func(t *testing.T) {
		// Create cart first
		userID := uuid.New()
		cartID := uuid.New()
		cart := &cartEntity.Cart{
			ID:     cartID,
			UserID: userID,
		}

		addItemReq := fiber.Map{
			"cart_id":    cartID.String(),
			"product_id": uuid.New().String(),
			"quantity":   2,
			"price":      50.0,
		}
		req, _ := json.Marshal(addItemReq)
		resp, err := app.Test(httptest.NewRequest("POST", "/api/cart/items", bytes.NewReader(req)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSagaAPI(t *testing.T) {
	app, tdb := setupTestApp(t)
	defer tdb.Cleanup()

	t.Run("start and monitor saga", func(t *testing.T) {
		// Create test order first
		orderID := uuid.New()
		order := &orderEntity.Order{
			ID:          orderID,
			UserID:      uuid.New(),
			TotalAmount: 100.0,
			Status:      orderEntity.OrderStatusPending,
		}

		// Start saga
		req := fiber.Map{
			"order_id": orderID.String(),
		}
		reqBody, _ := json.Marshal(req)
		resp, err := app.Test(httptest.NewRequest("POST", "/api/saga/order-payment", bytes.NewReader(reqBody)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check saga status
		resp, err = app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/saga/%s", orderID), nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var sagaStatus map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&sagaStatus)
		require.NoError(t, err)
		assert.Contains(t, []string{"PENDING", "COMPLETED", "FAILED"}, sagaStatus["status"])
	})
}
