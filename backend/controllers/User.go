package controllers

import (
	api "awesomeProject/openAPIGenerated"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"

	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// CreateUser handles creating a new user
func (uc *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Email        string `json:"email" binding:"required,email"`
		PasswordHash string `json:"passwordHash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ResourceID:   uuid.New(),
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	respUser := api.User{
		CreatedAt: &user.CreatedAt,
		Email:     openapi_types.Email(user.Email),
		Id:        user.ResourceID,
		UpdatedAt: &user.UpdatedAt,
	}
	c.JSON(http.StatusCreated, respUser)
}

// UpdateUser handles updating an existing user's email and password
func (uc *UserController) UpdateUser(c *gin.Context, id uuid.UUID) {
	var input struct {
		Email        string `json:"email" binding:"omitempty,email"`
		PasswordHash string `json:"passwordHash" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := uc.DB.Where("resource_id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.PasswordHash != "" {
		user.PasswordHash = input.PasswordHash
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
