package controllers

import (
	"awesomeProject/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" format:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479" doc:"The unique identifier of the user"`
	Email     string    `json:"email" format:"email" example:"user@example.com" doc:"The user's email address"`
	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2023-12-01T12:00:00Z" doc:"The time when the user was created"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time" example:"2023-12-02T15:00:00Z" doc:"The last time the user's details were updated"`
}

// RegisterUserInput represents the input for user registration
type RegisterUserInput struct {
	Body struct {
		Email    string `json:"email" format:"email" example:"user@example.com" doc:"User's email address"`
		Password string `json:"password" format:"password" example:"StrongPass!123" doc:"User's password"`
	}
}

// RegisterUserOutput represents the output for user registration
type RegisterUserOutput struct {
	Body User
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	ID   string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the user"`
	Body struct {
		Email    string `json:"email,omitempty" format:"email" example:"newemail@example.com" doc:"User's new email address"`
		Password string `json:"password,omitempty" format:"password" example:"NewPass123!" doc:"User's new password"`
	}
}

// UpdateUserOutput represents the output for updating a user
type UpdateUserOutput struct {
	Body User
}

type UserController struct {
	DB *gorm.DB
}

// CreateUser handles creating a new user
func (uc *UserController) CreateUser(ctx context.Context, input *RegisterUserInput) (*RegisterUserOutput, error) {
	user := models.User{
		ResourceID:   uuid.New(),
		Email:        input.Body.Email,
		PasswordHash: input.Body.Password,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		// Translate DB error to huma error HTTP code
		// TODO: Maybe make it as something like a post-processing middleware?
		return nil, ErrorGormToHuma(err)
	}
	resp := &RegisterUserOutput{}
	resp.Body.ID = user.ResourceID.String()
	resp.Body.Email = user.Email
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt
	return resp, nil
}

func (uc *UserController) UpdateUser(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
	var user models.User
	if err := uc.DB.Where("resource_id = ?", input.ID).First(&user).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	if input.Body.Email != "" {
		user.Email = input.Body.Email
	}
	if input.Body.Password != "" {
		user.PasswordHash = input.Body.Password
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}
	resp := &UpdateUserOutput{
		Body: User{
			ID:        user.ResourceID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
	return resp, nil
}
