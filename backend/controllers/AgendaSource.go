package controllers

import (
	"awesomeProject/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgendaSource represents an agenda source in the API
type AgendaSource struct {
	ID        string    `json:"id" format:"uuid" example:"c29ac10b-58cc-4372-a567-0e02b2c3d479" doc:"The unique identifier of the agenda source"`
	URL       string    `json:"url" format:"uri" example:"https://example.com/calendar" doc:"The URL of the agenda source"`
	Type      string    `json:"type" enum:"proton" example:"proton" doc:"The type of the agenda source"`
	UserID    string    `json:"userId" format:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479" doc:"The ID of the user who owns the agenda source"`
	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2023-12-01T12:00:00Z" doc:"The time when the agenda source was created"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time" example:"2023-12-02T15:00:00Z" doc:"The last time the agenda source was updated"`
}

// GetAgendaSourcesInput represents the input for getting agenda sources
type GetAgendaSourcesInput struct {
	OrderBy  string `query:"orderBy" enum:"asc,desc" default:"asc" doc:"Order the results by 'updatedAt' in ascending ('asc') or descending ('desc') order."`
	Page     int    `query:"page" minimum:"1" default:"1" doc:"The page number to retrieve (1-based)."`
	PageSize int    `query:"pageSize" minimum:"1" maximum:"100" default:"20" doc:"The number of items to include per page."`
}

// GetAgendaSourcesOutput represents the output for getting agenda sources
type GetAgendaSourcesOutput struct {
	Body struct {
		Data       []AgendaSource `json:"data"`
		Pagination Pagination     `json:"pagination"`
	}
}

// CreateAgendaSourceInput represents the input for creating an agenda source
type CreateAgendaSourceInput struct {
	Body struct {
		URL  string `json:"url" format:"uri" example:"https://example.com/calendar" doc:"The URL of the agenda source"`
		Type string `json:"type" enum:"proton" example:"proton" doc:"The type of the agenda source"`
	}
}

// CreateAgendaSourceOutput represents the output for creating an agenda source
type CreateAgendaSourceOutput struct {
	Body AgendaSource
}

// GetAgendaSourceInput represents the input for getting an agenda source
type GetAgendaSourceInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda source"`
}

// GetAgendaSourceOutput represents the output for getting an agenda source
type GetAgendaSourceOutput struct {
	Body AgendaSource
}

// UpdateAgendaSourceInput represents the input for updating an agenda source
type UpdateAgendaSourceInput struct {
	ID   string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda source"`
	Body struct {
		URL  string `json:"url,omitempty" format:"uri" example:"https://newexample.com/calendar" doc:"The URL of the agenda source"`
		Type string `json:"type,omitempty" enum:"proton" example:"proton" doc:"The type of the agenda source"`
	}
}

// UpdateAgendaSourceOutput represents the output for updating an agenda source
type UpdateAgendaSourceOutput struct {
	Body AgendaSource
}

// DeleteAgendaSourceInput represents the input for deleting an agenda source
type DeleteAgendaSourceInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda source"`
}

// AgendaSourceController handles operations on agenda sources
type AgendaSourceController struct {
	DB *gorm.DB
}

// GetAgendaSources retrieves a list of agenda sources with pagination
func (asc *AgendaSourceController) GetAgendaSources(ctx context.Context, input *GetAgendaSourcesInput) (*GetAgendaSourcesOutput, error) {
	var agendaSources []models.AgendaSource
	var count int64

	// Get user ID from context (assuming it's set by authentication middleware)
	// For now, we'll just query all sources

	// Set up pagination
	offset := (input.Page - 1) * input.PageSize

	// Set up ordering
	order := "updated_at"
	if input.OrderBy == "desc" {
		order = "updated_at DESC"
	}

	// Count total items
	if err := asc.DB.Model(&models.AgendaSource{}).Count(&count).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Get paginated items
	if err := asc.DB.Order(order).Offset(offset).Limit(input.PageSize).Find(&agendaSources).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Prepare response
	resp := &GetAgendaSourcesOutput{}
	resp.Body.Data = make([]AgendaSource, len(agendaSources))

	// Convert model to API response
	for i, source := range agendaSources {
		resp.Body.Data[i] = AgendaSource{
			ID:        source.ResourceID.String(),
			URL:       source.Url,
			Type:      source.Type,
			UserID:    uuid.UUID{}.String(), // This should be replaced with actual user ID
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		}
	}

	// Set pagination info
	totalPages := int(count) / input.PageSize
	if int(count)%input.PageSize > 0 {
		totalPages++
	}

	resp.Body.Pagination = Pagination{
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalItems: int(count),
		TotalPages: totalPages,
	}

	return resp, nil
}

// CreateAgendaSource creates a new agenda source
func (asc *AgendaSourceController) CreateAgendaSource(ctx context.Context, input *CreateAgendaSourceInput) (*CreateAgendaSourceOutput, error) {
	// Create a new agenda source
	agendaSource := models.AgendaSource{
		ResourceID: uuid.New(),
		Url:        input.Body.URL,
		Type:       input.Body.Type,
		// UserID should be set from authenticated user
		// For now, we'll just set it to 1
		UserID: 1,
	}

	// Save to database
	if err := asc.DB.Create(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Prepare response
	resp := &CreateAgendaSourceOutput{}
	resp.Body = AgendaSource{
		ID:        agendaSource.ResourceID.String(),
		URL:       agendaSource.Url,
		Type:      agendaSource.Type,
		UserID:    uuid.UUID{}.String(), // This should be replaced with actual user ID
		CreatedAt: agendaSource.CreatedAt,
		UpdatedAt: agendaSource.UpdatedAt,
	}

	return resp, nil
}

// GetAgendaSource retrieves a single agenda source by ID
func (asc *AgendaSourceController) GetAgendaSource(ctx context.Context, input *GetAgendaSourceInput) (*GetAgendaSourceOutput, error) {
	var agendaSource models.AgendaSource

	// Parse UUID from string
	resourceID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	// Find the agenda source
	if err := asc.DB.Where("resource_id = ?", resourceID).First(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Prepare response
	resp := &GetAgendaSourceOutput{}
	resp.Body = AgendaSource{
		ID:        agendaSource.ResourceID.String(),
		URL:       agendaSource.Url,
		Type:      agendaSource.Type,
		UserID:    uuid.UUID{}.String(), // This should be replaced with actual user ID
		CreatedAt: agendaSource.CreatedAt,
		UpdatedAt: agendaSource.UpdatedAt,
	}

	return resp, nil
}

// UpdateAgendaSource updates an existing agenda source
func (asc *AgendaSourceController) UpdateAgendaSource(ctx context.Context, input *UpdateAgendaSourceInput) (*UpdateAgendaSourceOutput, error) {
	var agendaSource models.AgendaSource

	// Parse UUID from string
	resourceID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	// Find the agenda source
	if err := asc.DB.Where("resource_id = ?", resourceID).First(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Update fields if provided
	if input.Body.URL != "" {
		agendaSource.Url = input.Body.URL
	}
	if input.Body.Type != "" {
		agendaSource.Type = input.Body.Type
	}

	// Save changes
	if err := asc.DB.Save(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Prepare response
	resp := &UpdateAgendaSourceOutput{}
	resp.Body = AgendaSource{
		ID:        agendaSource.ResourceID.String(),
		URL:       agendaSource.Url,
		Type:      agendaSource.Type,
		UserID:    uuid.UUID{}.String(), // This should be replaced with actual user ID
		CreatedAt: agendaSource.CreatedAt,
		UpdatedAt: agendaSource.UpdatedAt,
	}

	return resp, nil
}

// DeleteAgendaSource deletes an agenda source by ID
func (asc *AgendaSourceController) DeleteAgendaSource(ctx context.Context, input *DeleteAgendaSourceInput) (*struct{}, error) {
	var agendaSource models.AgendaSource

	// Parse UUID from string
	resourceID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	// Find the agenda source
	if err := asc.DB.Where("resource_id = ?", resourceID).First(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Delete the agenda source
	if err := asc.DB.Delete(&agendaSource).Error; err != nil {
		return nil, ErrorGormToHuma(err)
	}

	// Return empty response for 204 No Content
	return &struct{}{}, nil
}
