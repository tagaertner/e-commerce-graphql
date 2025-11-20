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
	"github.com/tagaertner/e-commerce-graphql/services/users/models"
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

    require.NoError(t, db.AutoMigrate(&models.User{}, &models.User{}))
    return db
}

// setupTestEnv initializes a fresh test environment for User service tests.
// It sets up the database, creates a new UserService instance, 
// and returns the DB, service, and context for use within tests.
func setupTestEnv(t *testing.T) (*gorm.DB, *UserService, context.Context) {
	db := setupTestDB(t)
	userService := NewUserService(db)
	ctx := context.Background()
	return db, userService, ctx
}

// Helper Function 
func BoolPointer(b bool) *bool{
	return &b
}

func StringPointer(s string) *string {
	return &s
}

// === Tests ===

// 🧪 GetUserByID
// 	1.	TestGetUserByID_ReturnsUser_WhenExists
func TestGetUserByID_ReturnsUser_WhenExists(t *testing.T){
	_, userService, ctx := setupTestEnv(t)

	// create user
	created, err := userService.CreateUser(
		ctx,
		"Tina Test",
		"tinatest@test.com",
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)

	// Retrieve that user by its generated id
	found, err := userService.GetUserByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)

	// Verify the data
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "Tina Test", found.Name)
	assert.Equal(t, "tinatest@test.com", found.Email)

}

// 	2.	TestGetUserByID_ReturnsError_WhenUserNotFound
func TestGetUserByID_ReturnsError_WhenUserNotFound(t *testing.T){
	_, userService, ctx := setupTestEnv(t)

	created, err := userService.GetUserByID(
		ctx,
		"",
		
	)
	assert.EqualError(t, err,"record not found")
	assert.Nil(t, created, "user should return error when userID is not found")

}

// 🧪 GetAllUsers
// 	3.	TestGetAllUsers_ReturnsAllUsers
func TestGetAllUsers_ReturnsAllUsers(t *testing.T){
	db, userService, ctx := setupTestEnv(t)
	
	// create user
	user1 := models.User{
	ID:     "user1",
	Name:   "1User",
	Email:  "1user@test.com",
	Role:   "user",
	Active: true,
	}
	user2 := models.User{
	ID:     "user2",
	Name:   "2User",
	Email:  "2user@test.com",
	Role:   "user",
	Active: true,
	}
	user3 := models.User{
	ID:     "user3",
	Name:   "3User",
	Email:  "3user@test.com",
	Role:   "user",
	Active: true,
	}

	require.NoError(t,db.Create(&user1).Error)
	require.NoError(t,db.Create(&user2).Error)
	require.NoError(t,db.Create(&user3).Error)

	// ---Act ---
	users, err := userService.GetAllUsers(
		ctx,
	
	)

	// --Assert ---
	require.NoError(t, err)
	require.Len(t, users,3, "should return all users")

	// Verify all 3 users exist
	names := []string{users[0].Name, users[1].Name, users[2].Name}
	assert.Contains(t, names, "1User")
	assert.Contains(t, names, "2User")
	assert.Contains(t, names, "3User")

}
// 	4.	TestGetAllUsers_ReturnsEmptySlice_WhenNoUsersExist
func TestGetAllUsers_ReturnsEmptySlice_WhenNoUsersExist(t *testing.T){
	_, userService, ctx := setupTestEnv(t)

	// ---Arrange---
	// No Users created 

	// ---Act ---
	users, err := userService.GetAllUsers(ctx)

	// ---Assert ---
	require.NoError(t, err, "expected no error whne there are no users")
	require.NotNil(t, users, "expected slice , not nil")
	require.Len(t, users, 0, "expected empty slice when no users exist")
}


// 🧪 CreateUser
// TestCreateUser_Success
func TestCreateUser_Success(t *testing.T){
	_, userSevice, ctx := setupTestEnv(t)

	created, err := userSevice.CreateUser(
		ctx,
		"user1",
		"1user@test.com",
	)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, created)

	assert.NotEmpty(t, created.ID, "Id should be auto-generated")
	assert.Equal(t, "user1", created.Name, "Name should match input")
	assert.Equal(t, "1user@test.com", created.Email, "Email should match input")

	assert.Equal(t, "user",created.Role, "Default role should be 'user'")
	assert.True(t, created.Active, "User should be active by default")

}
// 	6.	TestCreateUser_ReturnsError_WhenDatabaseFails
func TestCreateUser_ReturnsError_WhenDatabaseFails(t *testing.T){
	_, userService, ctx := setupTestEnv(t)

	created, err := userService.CreateUser(
	
		ctx,
		" ",
		"3user@test.com",
	)
	assert.Error(t, err, "invalid database")
	assert.Nil(t, created, "user should not be created when userID or user name is invalid")
}
// 	7.	TestCreateUser_SetsDefaultRoleAndActiveFields
func TestCreateUser_SetsDefaultRoleAndActiveFields(t *testing.T){
 //TOdo
}
func TestUpdateUser_UpdatesProvidedFields(t *testing.T){
	db, userService, ctx := setupTestEnv(t)

	// ---Arrange ---
	user := models.User{
		ID: "o1",
		Name: "Nancy Drew",
		Email: "nancyTest@email.com",
		Role: "customer",
		Active: true,
	}
	require.NoError(t, db.Create(&user).Error)

	// prepare for input update
	// newEmail := "nancy2Test@email.com "
	
	
	// ---Act ---
	input := &models.UpdateUserInput{
		ID:    "o1",
		Email: StringPointer("nancy2Test@email.com "),
	}

	update, err := userService.UpdateUser(ctx, input)

	// ---Assert ---
	require.NoError(t, err)
	require.NotNil(t, update)
	require.Equal(t, "nancy2Test@email.com ", update.Email)
}
// 	9.	TestUpdateUser_DoesNotUpdateWhenNoFieldsProvided
func TestUpdateUser_DoesNotUpdateWhenNoFieldsProvided (t *testing.T){
	db, userService, ctx:= setupTestEnv(t)

	// --Arrange ---
	original := models.User{
		ID: "123",
    	Name: "John",
    	Email: "john@test.com",
   		Role: "user",
    	Active: true,
	}
	require.NoError(t, db.Create(&original).Error)

	// Prepare empty update (no fields set)
	input := &models.UpdateUserInput{
		ID: "123", // only ID provided
		// Name: nil
		// Email: nil
		// Role: nil
		// Active: nil
	}

	//---Act ---
	updated, err := userService.UpdateUser(ctx, input)
	require.NoError(t, err, "expected no error when quering non-existaent user")
	require.NoError(t, err, "no error should occur for empty updates")
	require.NotNil(t, updated)

	// --- Assert ---
	// Pull from DB to be absolutely sure nothing changed
	var fetched models.User
	require.NoError(t, db.First(&fetched, "id = ?", "123").Error)

	assert.Equal(t, original.Name, fetched.Name)
	assert.Equal(t, original.Email, fetched.Email)
	assert.Equal(t, original.Role, fetched.Role)
	assert.Equal(t, original.Active, fetched.Active)

}
// 	10.	TestUpdateUser_ReturnsError_WhenUserDoesNotExist
// 	11.	TestUpdateUser_PartiallyUpdatesFieldsSuccessfully



// 🧪 DeleteUser
// 	12.	TestDeleteUser_ByID_SuccessfullyDeletesUser
// 	13.	TestDeleteUser_ByName_SuccessfullyDeletesUser
// 	14.	TestDeleteUser_ReturnsFalse_WhenNoUserFound
// 	15.	TestDeleteUser_ReturnsError_WhenDatabaseFails

