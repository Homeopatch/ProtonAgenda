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
	err = db.AutoMigrate(&models.User{}, &models.AgendaSource{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// setupAPI creates a test API with the controllers
func setupAPI(t *testing.T, db *gorm.DB) humatest.TestAPI {
	_, api := humatest.New(t)

	// Create controllers with the real DB
	userController := &controllers.UserController{
		DB: db,
	}

	agendaSourceController := &controllers.AgendaSourceController{
		DB: db,
	}

	// Register routes using the addRoutes function
	addRoutes(api, userController, agendaSourceController)

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

func TestAgendaSourceCRUD(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Setup API
	api := setupAPI(t, db)

	// Test creating an agenda source
	t.Run("Create agenda source", func(t *testing.T) {
		// Create a context with authentication
		ctx := context.Background()

		// Make a request to create an agenda source
		resp := api.PostCtx(ctx, "/api/agenda-sources", map[string]interface{}{
			"url":  "https://example.com/calendar",
			"type": "proton",
		})

		// Check response status code
		assert.Equal(t, http.StatusOK, resp.Code)

		// Parse response body
		var responseBody struct {
			ID        string    `json:"id"`
			URL       string    `json:"url"`
			Type      string    `json:"type"`
			UserID    string    `json:"userId"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
		err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// Verify response fields
		assert.NotEmpty(t, responseBody.ID)
		assert.Equal(t, "https://example.com/calendar", responseBody.URL)
		assert.Equal(t, "proton", responseBody.Type)
		assert.NotEmpty(t, responseBody.UserID)
		assert.False(t, responseBody.CreatedAt.IsZero())
		assert.False(t, responseBody.UpdatedAt.IsZero())

		// Store the ID for later tests
		agendaSourceID := responseBody.ID

		// Test getting the agenda source
		t.Run("Get agenda source", func(t *testing.T) {
			// Create a context with authentication
			ctx := context.Background()

			// Make a request to get the agenda source
			resp := api.GetCtx(ctx, "/api/agenda-sources/"+agendaSourceID)

			// Check response status code
			assert.Equal(t, http.StatusOK, resp.Code)

			// Parse response body
			var responseBody struct {
				ID        string    `json:"id"`
				URL       string    `json:"url"`
				Type      string    `json:"type"`
				UserID    string    `json:"userId"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Verify response fields
			assert.Equal(t, agendaSourceID, responseBody.ID)
			assert.Equal(t, "https://example.com/calendar", responseBody.URL)
			assert.Equal(t, "proton", responseBody.Type)
		})

		// Test updating the agenda source
		t.Run("Update agenda source", func(t *testing.T) {
			// Create a context with authentication
			ctx := context.Background()

			// Make a request to update the agenda source
			resp := api.PutCtx(ctx, "/api/agenda-sources/"+agendaSourceID, map[string]interface{}{
				"url": "https://updated-example.com/calendar",
			})

			// Check response status code
			assert.Equal(t, http.StatusOK, resp.Code)

			// Parse response body
			var responseBody struct {
				ID        string    `json:"id"`
				URL       string    `json:"url"`
				Type      string    `json:"type"`
				UserID    string    `json:"userId"`
				CreatedAt time.Time `json:"createdAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Verify response fields
			assert.Equal(t, agendaSourceID, responseBody.ID)
			assert.Equal(t, "https://updated-example.com/calendar", responseBody.URL)
			assert.Equal(t, "proton", responseBody.Type)
		})

		// Test getting all agenda sources
		t.Run("Get all agenda sources", func(t *testing.T) {
			// Create a context with authentication
			ctx := context.Background()

			// Make a request to get all agenda sources
			resp := api.GetCtx(ctx, "/api/agenda-sources")

			// Check response status code
			assert.Equal(t, http.StatusOK, resp.Code)

			// Parse response body
			var responseBody struct {
				Data []struct {
					ID        string    `json:"id"`
					URL       string    `json:"url"`
					Type      string    `json:"type"`
					UserID    string    `json:"userId"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"data"`
				Pagination struct {
					Page       int `json:"page"`
					PageSize   int `json:"pageSize"`
					TotalItems int `json:"totalItems"`
					TotalPages int `json:"totalPages"`
				} `json:"pagination"`
			}
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Verify response fields
			assert.GreaterOrEqual(t, len(responseBody.Data), 1)
			assert.Equal(t, 1, responseBody.Pagination.Page)
			assert.Equal(t, 20, responseBody.Pagination.PageSize)
			assert.GreaterOrEqual(t, responseBody.Pagination.TotalItems, 1)
			assert.GreaterOrEqual(t, responseBody.Pagination.TotalPages, 1)

			// Find our agenda source in the list
			found := false
			for _, source := range responseBody.Data {
				if source.ID == agendaSourceID {
					assert.Equal(t, "https://updated-example.com/calendar", source.URL)
					assert.Equal(t, "proton", source.Type)
					found = true
					break
				}
			}
			assert.True(t, found, "Created agenda source not found in the list")
		})

		// Test deleting the agenda source
		t.Run("Delete agenda source", func(t *testing.T) {
			// Create a context with authentication
			ctx := context.Background()

			// Make a request to delete the agenda source
			resp := api.DeleteCtx(ctx, "/api/agenda-sources/"+agendaSourceID)

			// Check response status code - should be 204 No Content
			assert.Equal(t, http.StatusNoContent, resp.Code)

			// Try to get the deleted agenda source
			getResp := api.GetCtx(ctx, "/api/agenda-sources/"+agendaSourceID)

			// Check response status code - should be 404 Not Found
			assert.Equal(t, http.StatusNotFound, getResp.Code)
		})
	})

	// Test agenda source not found
	t.Run("Agenda source not found", func(t *testing.T) {
		// Create a context with authentication
		ctx := context.Background()

		// Generate a random UUID that doesn't exist in the database
		nonExistentID := uuid.New().String()

		// Make a request to get a non-existent agenda source
		resp := api.GetCtx(ctx, "/api/agenda-sources/"+nonExistentID)

		// Check response status code - should be not found
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
