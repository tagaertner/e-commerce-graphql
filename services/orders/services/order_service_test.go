package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagaertner/e-commerce-graphql/services/orders/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a new in-memory SQLite database for unit testing.
// It runs AutoMigrate on the Order model to prepare the schema.
// This database is isolated to the current test and discarded when the test ends.
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect test database: %v", err)
    }
    db.AutoMigrate(&models.Order{})
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
// todo GEt
// TestGetOrdersByUserID_Success confirms that the service correctly
// retrieves all orders for a given user from the test database.
func TestGetOrderByUserId_success(t *testing.T){
	// _, orderService, ctx := setupTestEnv(t)
}

// TestGetOrdersByUserID_Failure checks that the service handles cases
// where a user has no existing orders and returns an empty result set.
func TestGetOrdersByUserID_Failure(t *testing.T) {

}

//todo update
// TestUpdateOrderStatus_Success verifies that the service updates an
// existing order’s status correctly and persists the change in the database.
func TestUpdateOrderStatus_Success(t *testing.T) { 
}

// TestUpdateOrderStatus_Failure ensures that attempting to update a
// non-existent order returns an error and does not modify any data.
func TestUpdateOrderStatus_Failure(t *testing.T) { 

 }


// todo Delet
// TestDeleteOrder_Success validates that a created order can be deleted
// and no longer exists in the test database after deletion.
func TestDeleteOrder_Success(t *testing.T) { 

 }


// TestDeleteOrder_Failure checks that deleting an order with an invalid ID
// returns an error and does not affect existing records.
func TestDeleteOrder_Failure(t *testing.T) { 

 }


//todo TestCreateOrder_ZeroQuantity: should fail validation

// Todo TestGetOrderbyUserID_MultipleUsers: ensure only that users order return

// Todo TestUpdateOrderStatus_NoChange: updating with same status should not error

// Todo TestDeleteOrder_Twice : second delet should error gracefully
// Example GraphQL E2E tests:
//•	✅ TestQueryOrders — queries all orders and checks response shape.
//•	✅ TestMutationCreateOrder — hits GraphQL mutation (optional).