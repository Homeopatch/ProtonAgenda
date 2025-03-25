package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	userController := &UserController{DB: db}

	router.POST("/api/users", userController.CreateUser)
	router.PUT("/api/users/:id", userController.UpdateUser)

	return router
}

// Setup PostgreSQL database for tests
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

	//err := dropTestDatabase(host, port, user, password, dbname)
	//if err != nil {
	//	return nil, err
	//}

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("DROP TABLE IF EXISTS test_db")

	// Auto-migrate all models
	err = db.AutoMigrate(&models.User{}, &models.AgendaInvite{}, &models.AgendaSource{}, &models.AgendaItem{}, &models.ProceduralAgenda{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestUserController(t *testing.T) {
	db, _ := setupTestDB()
	router := setupRouter(db)

	var createdResourceID string

	tests := []struct {
		name       string
		method     string
		url        string
		body       map[string]string
		statusCode int
		validate   func(t *testing.T, db *gorm.DB)
	}{
		{
			name:   "Create user successfully",
			method: "POST",
			url:    "/api/users",
			body: map[string]string{
				"email":        "test@example.com",
				"passwordHash": "securepassword",
			},
			statusCode: http.StatusCreated,
			validate: func(t *testing.T, db *gorm.DB) {
				var user models.User
				err := db.First(&user, "email = ?", "test@example.com").Error
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", user.Email)

				// Capture the ResourceID for use in subsequent tests
				createdResourceID = user.ResourceID.String()
				assert.NotEmpty(t, createdResourceID)
			},
		},
		{
			name:   "Update user successfully",
			method: "PUT",
			url:    "", // Placeholder, we'll set this dynamically
			body: map[string]string{
				"email":        "updated@example.com",
				"passwordHash": "updatedpassword",
			},
			statusCode: http.StatusOK,
			validate: func(t *testing.T, db *gorm.DB) {
				var user models.User
				err := db.First(&user, "email = ?", "updated@example.com").Error
				assert.NoError(t, err)
				assert.Equal(t, "updatedpassword", user.PasswordHash)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set the correct URL for the update test case
			if tc.name == "Update user successfully" {
				tc.url = "/api/users/" + createdResourceID
				assert.NotEmpty(t, createdResourceID, "ResourceID should not be empty before running this test")
			}

			// Marshal request body
			var payload []byte
			if tc.body != nil {
				payload, _ = json.Marshal(tc.body)
			}

			// Create and send the request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Validate the response
			assert.Equal(t, tc.statusCode, w.Code)

			// Perform custom validation
			if tc.validate != nil {
				tc.validate(t, db)
			}
		})
	}
}
