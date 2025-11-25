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
	"github.com/tagaertner/e-commerce-graphql/services/products/models"
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

    require.NoError(t, db.AutoMigrate(&models.Product{}, &models.Product{}))
    return db
}

// setupTestEnv initializes a fresh test environment for Order service tests.
// It sets up the in-memory database, creates a new OrderService instance, 
// and returns the DB, service, and context for use within tests.
func setupTestEnv(t *testing.T) (*gorm.DB, *ProductService, context.Context) {
	db := setupTestDB(t)
	productService := NewProductService(db)
    ctx := context.Background()
	return db, productService, ctx
}

// Helper function for string literals
func strPtr(s string) *string {
    return &s
}

// === Tests ===

// TestCreateProduct_Success verifies that a product with valid name, price, inventory, and description
// is created successfully without errors.
func TestCreateProduct_Success(t *testing.T){
	_, productService, ctx := setupTestEnv(t)
	desc := "Simple widget"
	created, err := productService.CreateProduct(ctx, "Widget", 29.99, desc, 50)

	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Widget", created.Name)
	assert.Equal(t, float64(29.99), created.Price)
}

// TestCreateProduct_Failure verifies that inventory is updated
func TestCreateProduct_Failure(t *testing.T) {
	_, productService, ctx := setupTestEnv(t)
	desc := "Simple widget"
	created, err := productService.CreateProduct(
		ctx,
		"",
		29.99,
		desc,
		50,
	)
	assert.EqualError(t, err, "invalid product name: missing or invalid field")
	assert.Nil(t, created, "Product should not be created")
}

// 	TestUpdateProduct_Success update inventory to 8
func TestUpdateProduct_Success(t *testing.T){
	db, productService, ctx := setupTestEnv(t)

	// ---Arrange ---
	productDesc := "New Widget"
	product := models.Product{
		Name: "Widgets",
		Price: 4.99,
		Description: &productDesc,
		Inventory: 10,
	}
	require.NoError(t, db.Create(&product).Error)

	// Prepare input for update
	newInventory := 8
	updateInput := &models.UpdateProductInput{
		Inventory: &newInventory,
	}

	// ---Act---
	updated, err := productService.UpdateProduct(ctx, product.ID, *updateInput)

	// ---Assert---
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newInventory, updated.Inventory)
}


//TestUpdateProduct_Failure Update non-existent ID should error.
func TestUpdateProduct_Failure(t *testing.T){
	_, productService, ctx := setupTestEnv(t)

	// ---Arrange ---
	newName := " Update Widget"
	updateInput := &models.UpdateProductInput{
		Name: &newName,
	}

	// ---Act---
	updated, err := productService.UpdateProduct(ctx, "non-existent-id", *updateInput)

	// ---Assert ---
	assert.Error(t, err, "expected error when updating non-existent product")
	assert.Nil(t, updated, "no product should be returned on failure")
}

//TestDeleteProduct_Success Creates product, deletes it, then verifies it no longer exists.
func TestDeleteProduct_Success (t *testing.T){
	db, productService, ctx := setupTestEnv(t)

	// ---Arrange --- ID, Name
	product := models.Product{
		ID: "ab23",
		Name: "Old Widget",
	}
	
	require.NoError(t, db.Create(&product).Error)

	_, err := productService.DeleteProduct(ctx, models.DeleteProductInput{
		ID: &product.ID,
		Name: strPtr("Old Widget"),
	})
	require.NoError(t, err)

	var found models.Product
	result := db.First(&found, "id =?", product.ID)
	assert.Error(t, result.Error, "product should be deleted")

}

//TestDeleteProduct_Failure Delete with bad ID returns error, no rows affected.
func TestDeleteProduct_Failure(t *testing.T){
	db, productServices, ctx := setupTestEnv(t)

	var countBefore int64
	db.Model(&models.Product{}).Count(&countBefore)

	// Capture both return values and assert on the error only
	_, err := productServices.DeleteProduct(ctx, models.DeleteProductInput{
		ID: strPtr("non-existent-id"),
		Name: strPtr("Old Widget"),
	})
	assert.Error(t, err, "should return error for non-existant id")

	var countAfter int64
	db.Model(&models.Product{}).Count(&countAfter)
	assert.Equal(t, countBefore, countAfter, "no rows should be deleted")

}

// TestCreateProduct_ZeroOrNegativeInventory Reject inventory < 0.
func TestCreateProduct_ZeroOrNegativeInventory(t *testing.T){
	_, productService, ctx := setupTestEnv(t)

	// ---Arrange --- name, price, description, inventory
	created, err := productService.CreateProduct(
		ctx,
		"Zero Widget",
		51.99,
		"Something invisable but awesome",
		0,
		
	)

	// ---Assert---
	assert.Error(t, err, "should reject zero inventory")
	assert.Nil(t, created, "product should not be created with zero inventory")

	// Test negative inventory
	created, err = productService.CreateProduct(
		ctx,
		"Nil Widget",
		 29.99,
		"Not so awesome",
		-5,
	
	)

	assert.Error(t, err, "should reject negative inventory")
	assert.Nil(t, created, "product should not be created with negative inventory")
}

// TestDeleteProduct_Twice First delete succeeds; second delete returns “not found”. id and name
func TestDeleteProduct_Twice (t *testing.T){
	db, productService, ctx := setupTestEnv(t)

	// ---Arrange --
	product := models.Product{
		ID: "p2",
		Name: "Gadget",
	}
	require.NoError(t, db.Create(&product).Error)

	// ---Act: First dlection should succeed ---
	success, err := productService.DeleteProduct(ctx, models.DeleteProductInput{
		ID: &product.ID,
		Name: strPtr("Gadget"),
	})
	require.NoError(t, err, "first deletion should succeed")
	require.True(t, success, "first deletion should return true")

	// ---Act: Second deletion should fail gracefully ---
	success, err = productService.DeleteProduct(ctx, models.DeleteProductInput{
		ID: &product.ID,
		Name: strPtr("Gadget"),
	})

	// ---Assert ---
	assert.Error(t, err, "second deletion should return an error")
	assert.False(t, success, "second deletion should return false")
	assert.Contains(t, err.Error(), "not found", "error shoul dindicate product not found")
}

//TestGetAllProducts_Success Insert multiple products, call GetAllProducts(), verify correct count + order.
func TestGetAllProducts_Success(t *testing.T){
	db, productService, _ := setupTestEnv(t)
		// Create some products
	product1 := models.Product {
		
		ID:		"p1",
		Name:  	"Widget",
		Price: 	59.99,
		Description: strPtr("Fancy widget"),
		Inventory: 40,
		Available: true,
	}

	product2 := models.Product {
		ID:		"p2",
		Name:  	"Gadget",
		Price: 	38.99,
		Description: strPtr("Plain gadget"),
		Inventory: 51,
		Available: true,
	}

		product3 := models.Product {
		ID:		"p3",
		Name:  	"Widgy Gadget",
		Price: 	55.99,
		Description: strPtr("Fancy smancy gadget"),
		Inventory: 30,
		Available: true,
	}
	require.NoError(t, db.Create(&product1).Error)
	require.NoError(t,db.Create(&product2).Error)
	require.NoError(t,db.Create(&product3).Error)
	
	// ---Act ---
	products, err :=productService.GetAllProducts()

	// -- Assert ---
	require.NoError(t, err)
	require.Len(t, products, 3, "should return all inserted products")

	// Verfiy all 3 names exits
	names := []string{products[0].Name, products[1].Name, products[2].Name}
	assert.Contains(t, names, "Widget")
	assert.Contains(t, names, "Gadget")
	assert.Contains(t, names, "Widgy Gadget")
}
// TestGetProductByID_Success Valid ID returns correct product.
func TestGetProductByID_Success(t *testing.T){
	db, productService, _ := setupTestEnv(t)
		// Create some products
	product1 := models.Product {
		
		ID:		"p1",
		Name:  	"Widget",
		Price: 	59.99,
		Description: strPtr("Fancy widget"),
		Inventory: 40,
		Available: true,
	}

	product2 := models.Product {
		ID:		"p2",
		Name:  	"Gadget",
		Price: 	38.99,
		Description: strPtr("Plain gadget"),
		Inventory: 51,
		Available: true,
	}

		product3 := models.Product {
		ID:		"p3",
		Name:  	"Thingamajig",
		Price: 	19.99,
		Description: strPtr("Fancy smancy gadget"),
		Inventory: 30,
		Available: true,
	}
	require.NoError(t, db.Create(&product1).Error)
	require.NoError(t, db.Create(&product2).Error)
	require.NoError(t, db.Create(&product3).Error)

	// ---Act ---
	products, err :=productService.GetProductByID("p3")

	// ---Assert ---
	require.NoError(t, err)
	require.NotNil(t, products)
	assert.Equal(t, "p3", product3.ID)
	assert.Equal(t, "Thingamajig", product3.Name)
	assert.Equal(t, 19.99, product3.Price)
	
}

// TestGetProductByID_FailureInvalid ID returns nil and error.
func TestGetProductByID_Failure (t *testing.T){
	_, productService, _ := setupTestEnv(t)

	// ---Arrange ---
	// No products created for this

	// ---Act ---
	products, err := productService.GetProductByID("p2")

	// ---Assert---
	require.Error(t, err, "expected error when querying non-existent product")
	require.Nil(t, products, "expected nil product for non-existent ID")
	
}

//TestRestockProduct_Success 
func TestRestockProduct_Success(t *testing.T){
	db, productServices, ctx := setupTestEnv(t)

	// ---Arrange ---
	product := models.Product{
		ID:  "p1",
		Name: "Widget",
		Price: 9.99,
		Description: strPtr("Simple widget"),
		Inventory: 5,
		Available: true,
	}
	require.NoError(t, db.Create(&product).Error)

	// ---Act ---
	updated, err := productServices.RestockProduct(ctx, "p1", 10)

	// ---Assert ---
	require.NoError(t, err)
	assert.Equal(t, 15, updated.Inventory)
}

// TestRestockProduct_Failure
func TestRestockProduct_Failure(t *testing.T){
	_, productService, ctx := setupTestEnv(t)

	//---Act---
	product, err := productService.RestockProduct(ctx, "non-existent-id", 10)

	// ----Assert ---
	assert.Error(t, err, "should return error when restocking non-existent product")
	assert.Nil(t, product, "no product should be returned on failure")

}

// TestRestockProduct_Failure_NegativeAmount reject negative amount
func TestRestockProduct_Failure_NegativeAmount(t *testing.T){
	db, productService, ctx := setupTestEnv(t)

	product := models.Product {
		ID:		"p3",
		Name:  	"Thingamajig",
		Price: 	19.99,
		Description: strPtr("Fancy smancy gadget"),
		Inventory: 30,
		Available: true,
	}
	require.NoError(t, db.Create(&product).Error)

	updated, err := productService.RestockProduct(ctx, "p3", -10)

	// ---Assert---
	assert.Error(t, err, "should reject negative restock amount")
	assert.Nil(t, updated)
	
}

//TestSetProductAvailability
func TestSetProductAvailability_Success(t *testing.T){
	db, productService, ctx := setupTestEnv(t)

	product := models.Product{
		ID:          "p2",
		Name:        "Gadget",
		Price:       19.99,
		Description: strPtr("Basic gadget"),
		Inventory:   12,
		Available:   true,
	}
	require.NoError(t, db.Create(&product).Error)

	// ---Act ---
	updated , err := productService.SetProductAvailability(ctx, "p2", false)

	// --Assert --
	require.NoError(t, err)
	assert.False(t, updated.Available)
}

func TestSetProductAvailability_Failure(t *testing.T){
	_, productService, ctx := setupTestEnv(t)

	// ---Act ---
	product, err := productService.SetProductAvailability(ctx, "non-existent-id", false)

	// ---Assert ---
	assert.Error(t, err, "should return error when toggling non-existent product")
	assert.Nil(t, product)
}

// TestSetProductAvailability_Failure_NotFound
func TestSetProductAvailability_Failure_NotFound(t *testing.T){
	_,productService, ctx := setupTestEnv(t)

	// ---Act ---
	product, err := productService.SetProductAvailability(ctx, "non-existent-id", false)

	// --- Assert ---
	assert.Error(t, err, "should return error when product does not exist")
	assert.Nil(t, product, "expected nil product on failure")
}