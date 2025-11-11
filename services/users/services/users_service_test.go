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
// 	5.	TestCreateUser_SuccessfullyCreatesUser
// 	6.	TestCreateUser_ReturnsError_WhenDatabaseFails
// 	7.	TestCreateUser_SetsDefaultRoleAndActiveFields



// 🧪 UpdateUser
// 	8.	TestUpdateUser_UpdatesProvidedFields
// 	9.	TestUpdateUser_DoesNotUpdateWhenNoFieldsProvided
// 	10.	TestUpdateUser_ReturnsError_WhenUserDoesNotExist
// 	11.	TestUpdateUser_PartiallyUpdatesFieldsSuccessfully



// 🧪 DeleteUser
// 	12.	TestDeleteUser_ByID_SuccessfullyDeletesUser
// 	13.	TestDeleteUser_ByName_SuccessfullyDeletesUser
// 	14.	TestDeleteUser_ReturnsFalse_WhenNoUserFound
// 	15.	TestDeleteUser_ReturnsError_WhenDatabaseFails

