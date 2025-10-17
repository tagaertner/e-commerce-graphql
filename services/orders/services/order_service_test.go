package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tagaertner/e-commerce-graphql/services/orders/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// === SET UP ===
// setupTestDB starts a temporary Postgres container using testcontainers-go.
// Requires Docker to be running. The container is created automatically for tests
// and removed after they complete, providing an isolated Postgres instance that
// matches production behavior.
func setupTestDB(t *testing.T) *gorm.DB {
    ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
		}).WithStartupTimeout(60 * time.Second), // ⏳ give it a full minute
	}
	
    pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)

    host, _ := pgContainer.Host(ctx)
    port, _ := pgContainer.MappedPort(ctx, "5432/tcp")

    dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    require.NoError(t, err)

    require.NoError(t, db.AutoMigrate(&models.Order{}, &models.Product{}))
    return db
}

// setupTestEnv initializes a fresh test environment for Order service tests.
// It sets up the in-memory database, creates a new OrderService instance, 
// and returns the DB, service, and context for use within tests.
func setupTestEnv(t *testing.T) (*gorm.DB, *OrderService, context.Context) {
	db := setupTestDB(t)
    orderService := NewOrderService(db)
    ctx := context.Background()
	return db, orderService, ctx
}

// === Tests ===

// TestCreateOrder_Success verifies that a valid order can be created successfully
func TestCreateOrder_Success(t *testing.T) {
	_, orderService, ctx := setupTestEnv(t)
    created, err := orderService.CreateOrder(
        ctx,
        "1",                 // userID
        []string{"101"},     // productIDs
        2,                   // quantity
        49.99,               // totalPrice
        "pending",           // status
		// todo currently using time_scalar for time, could switch to simulate job-stories
        time.Now(),          // createdAt
    )

    assert.NoError(t, err)
    assert.NotEmpty(t, created.ID)
    assert.Equal(t, "1", created.UserID)
    assert.Equal(t, float64(49.99), created.TotalPrice)
    assert.Equal(t, "pending", created.Status)
}

// TestCreateOrder_Failure ensures that CreateOrder returns an error
// and does not persist data when required fields are missing or invalid.
func TestCreateOrder_Failure(t *testing.T){
	_, orderService, ctx := setupTestEnv(t)

    created, err := orderService.CreateOrder(
        ctx,
        "",                 // userID
        []string{""},     // productIDs
        2,                   // quantity
        49.99,               // totalPrice
        "pending",           // status
        time.Now(),          // createdAt
    )

    assert.EqualError(t, err, "invalid order input: missing or invalid fields")
    assert.Nil(t, created, "order should not be created when userID or productIDs are invalid")
   

}

// TestGetOrdersByUserID_Success confirms that the service correctly
// retrieves all orders for a given user from the test database.
func TestGetOrderByUserId_success(t *testing.T){
	db, orderService, _  := setupTestEnv(t)
	
	// Create some products
	product1 := models.Product {
		
		ID:		"p1",
		Name:  	"Widget",
		Price: 	59.99,
		Description: "Fancy widget",
		Inventory: 40,
		Available: true,
	}

	product2 := models.Product {
		ID:		"p2",
		Name:  	"Gadget",
		Price: 	38.99,
		Description: "Plain gadget",
		Inventory: 51,
		Available: true,
	}
	require.NoError(t, db.Create(&product1).Error)
	require.NoError(t,db.Create(&product2).Error)

	// Create two orders for user 1
	order1 :=models.Order{
		ID:        "o1", 
		UserID:    "1",  
		Products: []models.Product{product1},
		Quantity:   1,    
		TotalPrice: 59.99,
		Status :    "pending",
		CreatedAt:	models.Now(),  
	}

	order2 :=models.Order{
		ID:        "02", 
		UserID:    "1",  
		Products: []models.Product{product2},
		Quantity:   1,    
		TotalPrice: 38.99,
		Status :    "completed",
		CreatedAt:	models.Now(),  
	}

	require.NoError(t, db.Create(&order1).Error)
	require.NoError(t, db.Create(&order2).Error)

	// --- Act ---
	orders, err := orderService.GetOrdersByUserID("1")

	// --- Assert ---
	require.NoError(t, err)
	require.Len(t, orders, 2)

	foundWidget, foundGadget := false, false

	for _, o := range orders {
		require.Len(t, o.Products, 1)
		switch o.Products[0].Name {
		case "Widget":
			require.Equal(t, "pending", o.Status)
			foundWidget = true
		case "Gadget":
			require.Equal(t, "completed", o.Status)
			foundGadget = true
		}
	}

	require.True(t, foundWidget, "Widget order not found")
	require.True(t, foundGadget, "Gadget order not found")
	
	// Validate preload worked
	require.Len(t, orders[0].Products,1)
	require.Equal(t, "Widget", orders[0].Products[0].Name)
	require.Equal(t, "pending", orders[0].Status)

	require.Len(t, orders[1].Products, 1)
	require.Equal(t, "Gadget", orders[1].Products[0].Name)
	require.Equal(t, "completed", orders[1].Status)

}

// TestGetOrdersByUserID_Failure checks that the service handles cases
// where a user has no existing orders and returns an empty result set.
func TestGetOrdersByUserID_Failure(t *testing.T) {
	_, orderService, _ := setupTestEnv(t)

	// --- Arrange ---
	// No orders created for this user ("999")

	// --- Act ---
	orders, err := orderService.GetOrdersByUserID("999")

	// -- Assert ---
	require.NoError(t, err, "expected no error when querying non-existent user")
	require.NotNil(t, orders, "expected a valid (non-nil) slice")
	require.Len(t, orders, 0, "expected no orders for this user")
}

//todo update
// TestUpdateOrderStatus_Success verifies that the service updates an
// existing order’s status correctly and persists the change in the database.
func TestUpdateOrderStatus_Success(t *testing.T) { 
	db, orderService, ctx := setupTestEnv(t)

	// --- Arrange ---
	order := models.Order{
		ID:         "o1",
		UserID:     "1",
		TotalPrice: 10.00,
		Status:     "pending",
	}
	require.NoError(t, db.Create(&order).Error)

	// Prepare input for update
	newStatus := "shipped"
	input := &models.UpdateOrderInput{
		OrderID: order.ID,
		Status:  &newStatus, // field is *string in your model
	}

	// --- Act ---
	updated, err := orderService.UpdateOrder(ctx, input)

	// --- Assert ---
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "shipped", updated.Status)
}


// TestUpdateOrderStatus_Failure ensures that attempting to update a
// non-existent order returns an error and does not modify any data.
func TestUpdateOrderStatus_Failure(t *testing.T) { 
	_, orderService, ctx := setupTestEnv(t)

	// --- Arrange ---
	// Prepare an input for a non-existent order ID
	newStatus := "shipped"
	input := &models.UpdateOrderInput{
		OrderID: "bad_id", // does not exist in DB
		Status:  &newStatus,
	}

	// --- Act ---
	updated, err := orderService.UpdateOrder(ctx, input)

	// --- Assert ---
	assert.Error(t, err, "expected an error when updating a non-existent order")
	assert.Nil(t, updated, "expected no order returned when update fails")
 }


// todo Delete
// TestDeleteOrder_Success validates that a created order can be deleted
// and no longer exists in the test database after deletion.
func TestDeleteOrder_Success(t *testing.T) { 
	db, orderService, ctx := setupTestEnv(t)

	// --- Arrange ---
	order := models.Order{
		ID:         "o1",
		UserID:     "1",
		TotalPrice: 10.00,
		Status:     "pending",
	}
	require.NoError(t, db.Create(&order).Error)

	// Capture both return values from DeleteOrder and use the created order's ID
	_, err := orderService.DeleteOrder(ctx, models.DeleteOrderInput{
		OrderID: order.ID,
		UserID: "1",
	})
	require.NoError(t, err)

	var found models.Order
	result := db.First(&found, "id = ?", order.ID)
	assert.Error(t, result.Error, "record should be deleted")
}

// TestDeleteOrder_Failure checks that deleting an order with an invalid ID
// returns an error and does not affect existing records.
func TestDeleteOrder_Failure(t *testing.T) { 
	db, orderService, ctx := setupTestEnv(t)
	
	var countBefore int64
	db.Model(&models.Order{}).Count(&countBefore)

	// Capture both return values and assert on the error only
	_, err := orderService.DeleteOrder(ctx, models.DeleteOrderInput{
		OrderID: "non-existent-id",
		UserID: "1",
	})
	assert.Error(t, err, "should return error for non-existent order")

	var countAfter int64
	db.Model(&models.Order{}).Count(&countAfter)
	assert.Equal(t, countBefore, countAfter, "no rows should be deleted")
}

//todo TestCreateOrder_ZeroQuantity: should fail validation

// Todo TestGetOrderbyUserID_MultipleUsers: ensure only that users order return

// Todo TestUpdateOrderStatus_NoChange: updating with same status should not error

// Todo TestDeleteOrder_Twice : second delet should error gracefully
// Example GraphQL E2E tests:
//•	✅ TestQueryOrders — queries all orders and checks response shape.
//•	✅ TestMutationCreateOrder — hits GraphQL mutation (optional).