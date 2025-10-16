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

// 
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test databaseL %v", err)
	}
	db.AutoMigrate(&models.Order{})
	return db
}

// Create
func TestCreateOrder_Success(t *testing.T) {
    db := setupTestDB(t)
    orderService := NewOrderService(db)

    ctx := context.Background()

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

// Failure Missing filed
func TestCreateOrder_Failure(t *testing.T){
	db := setupTestDB(t)
    orderService := NewOrderService(db)

    ctx := context.Background()

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
// Getorderby userID success
func TestGetOrderByUserId_success(t *testing.T){

}
// orderbyUserId failure noResutls

//todo update
// orderStatus success
// orderstatus fail invalid id


// todo Delet
// Order sucess
// order failer invaid id 




// Example GraphQL E2E tests:
//•	✅ TestQueryOrders — queries all orders and checks response shape.
//•	✅ TestMutationCreateOrder — hits GraphQL mutation (optional).