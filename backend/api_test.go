package main

import (
	"awesomeProject/controllers"
	"awesomeProject/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a test database connection
func setupTestDB() (*gorm.DB, error) {
	// Use environment variables or default values
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "test_db"
	}

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// setupAPI creates a test API with the user controller
func setupAPI(t *testing.T, db *gorm.DB) humatest.TestAPI {
	_, api := humatest.New(t)

	// Create a user controller with the real DB
	userController := &controllers.UserController{
		DB: db,
	}

	// Register routes using the addRoutes function
	addRoutes(api, userController)

	return api
}

func TestCreateUser(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Setup API
	api := setupAPI(t, db)

	// Test successful user creation
	t.Run("Successful user registration", func(t *testing.T) {
		// Make a request to register a user
		resp := api.Post("/api/register", map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		})

		// Check response status code
		assert.Equal(t, http.StatusOK, resp.Code)

		// Parse response body
		var responseBody struct {
			ID        string    `json:"id"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
		err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// Verify response fields
		assert.NotEmpty(t, responseBody.ID)
		assert.Equal(t, "test@example.com", responseBody.Email)
		assert.False(t, responseBody.CreatedAt.IsZero())
		assert.False(t, responseBody.UpdatedAt.IsZero())
	})
}

func TestUpdateUser(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Setup API
	api := setupAPI(t, db)

	// First, create a user through the API
	createResp := api.Post("/api/register", map[string]interface{}{
		"email":    "update-test@example.com",
		"password": "password123",
	})

	// Check that user was created successfully
	assert.Equal(t, http.StatusOK, createResp.Code)

	// Parse the response to get the user ID
	var createResponseBody struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
	err = json.Unmarshal(createResp.Body.Bytes(), &createResponseBody)
	assert.NoError(t, err)
	userID := createResponseBody.ID

	// Test successful user update
	t.Run("Successful user update", func(t *testing.T) {
		// Create a context with authentication
		ctx := context.Background()

		// Make a request to update a user
		resp := api.PutCtx(ctx, "/api/users/"+userID, map[string]interface{}{
			"email":    "updated@example.com",
			"password": "newpassword123",
		})

		// Check response status code
		assert.Equal(t, http.StatusOK, resp.Code)

		// Parse response body
		var responseBody struct {
			ID        string    `json:"id"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
		err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// Verify response fields
		assert.Equal(t, userID, responseBody.ID)
		assert.Equal(t, "updated@example.com", responseBody.Email)
		assert.False(t, responseBody.CreatedAt.IsZero())
		assert.False(t, responseBody.UpdatedAt.IsZero())
	})

	// Test user not found
	t.Run("User not found", func(t *testing.T) {
		// Create a context with authentication
		ctx := context.Background()

		// Generate a random UUID that doesn't exist in the database
		nonExistentID := uuid.New().String()

		// Make a request to update a non-existent user
		resp := api.PutCtx(ctx, "/api/users/"+nonExistentID, map[string]interface{}{
			"email": "updated@example.com",
		})

		// Check response status code - should be not found
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
