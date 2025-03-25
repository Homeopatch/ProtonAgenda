package main

import (
	"awesomeProject/controllers"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI
type Options struct {
	DbHost string `help:"Database hostname" env:"POSTGRES_HOST" default:"localhost"`
	DbPort int    `help:"Database port" env:"POSTGRES_PORT" default:"5432"`
	DbName string `help:"Database name" env:"POSTGRES_DBNAME" default:"test_db"`
	DbUser string `help:"Database username" env:"POSTGRES_USER" default:"postgres"`
	DbPass string `help:"Database password" env:"POSTGRES_PASSWORD" default:"password"`
	Port   int    `help:"Port to listen on" short:"p" default:"8888"`
}

// AgendaSource represents an agenda source to fetch agenda items from
type AgendaSource struct {
	ID        string    `json:"id" format:"uuid" example:"c29ac10b-58cc-4372-a567-0e02b2c3d479" doc:"The unique identifier of the agenda source"`
	URL       string    `json:"url" format:"uri" example:"https://example.com/calendar" doc:"The URL of the agenda source"`
	Type      string    `json:"type" enum:"proton" example:"proton" doc:"The type of the agenda source"`
	UserID    string    `json:"userId" format:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479" doc:"The ID of the user who owns the agenda source"`
	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2023-12-01T12:00:00Z" doc:"The time when the agenda source was created"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time" example:"2023-12-02T15:00:00Z" doc:"The last time the agenda source was updated"`
}

// AgendaItem represents an agenda item
type AgendaItem struct {
	ResourceID     string    `json:"ResourceID" format:"uuid" doc:"The unique identifier of the agenda item"`
	StartTime      time.Time `json:"StartTime" format:"date-time"`
	EndTime        time.Time `json:"EndTime" format:"date-time"`
	Description    string    `json:"Description"`
	AgendaSourceID string    `json:"AgendaSourceID" format:"uuid"`
	UserID         string    `json:"UserID" format:"uuid"`
}

// AgendaInvite represents an invitation to view a user's agenda
type AgendaInvite struct {
	ResourceID    string         `json:"ResourceID" format:"uuid" doc:"The unique identifier of the agenda invite"`
	UserID        string         `json:"UserID" format:"uuid" doc:"The ID of the user associated with the invite"`
	Description   string         `json:"Description"`
	ExpiresAt     time.Time      `json:"ExpiresAt" format:"date-time"`
	NotBefore     time.Time      `json:"NotBefore" format:"date-time"`
	NotAfter      time.Time      `json:"NotAfter" format:"date-time"`
	PaddingBefore string         `json:"PaddingBefore" doc:"Duration before the event"`
	PaddingAfter  string         `json:"PaddingAfter" doc:"Duration after the event"`
	SlotSizes     []string       `json:"SlotSizes" doc:"Array of slot sizes as durations"`
	AgendaSources []AgendaSource `json:"AgendaSources"`
}

// AgendaItemView represents a view of an agenda item without sensitive user data
type AgendaItemView struct {
	StartTime   time.Time `json:"StartTime" format:"date-time"`
	EndTime     time.Time `json:"EndTime" format:"date-time"`
	Description string    `json:"Description"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page" doc:"The current page number." example:"1"`
	PageSize   int `json:"pageSize" doc:"The number of items per page." example:"20"`
	TotalItems int `json:"totalItems" doc:"The total number of items available." example:"123"`
	TotalPages int `json:"totalPages" doc:"The total number of pages available." example:"7"`
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

// CreateAgendaItemsInput represents the input for creating agenda items
type CreateAgendaItemsInput struct {
	Body []AgendaItem
}

// GetAgendaItemsInput represents the input for getting agenda items
type GetAgendaItemsInput struct {
	AgendaSourceID string    `query:"agendaSourceID,omitempty" format:"uuid"`
	UserID         string    `query:"userID,omitempty" format:"uuid"`
	StartTime      time.Time `query:"startTime,omitempty" format:"date-time"`
	EndTime        time.Time `query:"endTime,omitempty" format:"date-time"`
	Page           int       `query:"page" minimum:"1" default:"1"`
	PageSize       int       `query:"pageSize" minimum:"1" maximum:"100" default:"20"`
}

// GetAgendaItemsOutput represents the output for getting agenda items
type GetAgendaItemsOutput struct {
	Body struct {
		Data       []AgendaItem `json:"data"`
		Pagination Pagination   `json:"pagination"`
	}
}

// GetAgendaItemInput represents the input for getting an agenda item
type GetAgendaItemInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda item"`
}

// GetAgendaItemOutput represents the output for getting an agenda item
type GetAgendaItemOutput struct {
	Body AgendaItem
}

// DeleteAgendaItemInput represents the input for deleting an agenda item
type DeleteAgendaItemInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda item"`
}

// CreateAgendaInviteInput represents the input for creating an agenda invite
type CreateAgendaInviteInput struct {
	Body AgendaInvite
}

// CreateAgendaInviteOutput represents the output for creating an agenda invite
type CreateAgendaInviteOutput struct {
	Body AgendaInvite
}

// GetAgendaInviteInput represents the input for getting an agenda invite
type GetAgendaInviteInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda invite"`
}

// GetAgendaInviteOutput represents the output for getting an agenda invite
type GetAgendaInviteOutput struct {
	Body AgendaInvite
}

// UpdateAgendaInviteInput represents the input for updating an agenda invite
type UpdateAgendaInviteInput struct {
	ID   string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda invite"`
	Body AgendaInvite
}

// UpdateAgendaInviteOutput represents the output for updating an agenda invite
type UpdateAgendaInviteOutput struct {
	Body AgendaInvite
}

// DeleteAgendaInviteInput represents the input for deleting an agenda invite
type DeleteAgendaInviteInput struct {
	ID string `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda invite"`
}

// ViewAgendaInviteInput represents the input for viewing an agenda invite
type ViewAgendaInviteInput struct {
	ID       string    `path:"id" format:"uuid" doc:"The unique identifier (UUID) of the agenda invite"`
	DateFrom time.Time `query:"DateFrom,omitempty" format:"date-time" doc:"The start date and time for filtering agenda items"`
	DateTo   time.Time `query:"DateTo,omitempty" format:"date-time" doc:"The end date and time for filtering agenda items"`
}

// ViewAgendaInviteOutput represents the output for viewing an agenda invite
type ViewAgendaInviteOutput struct {
	Body []AgendaItemView
}

func main() {
	// Create a CLI app which takes a port option
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()

		config := huma.DefaultConfig("User and Agenda Source Management API", "1.1.0")
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{

			// Example alternative describing the use of JWTs without documenting how
			// they are issued or which flows might be supported. This is simpler but
			// tells clients less information.
			"BearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		}
		api := humachi.New(router, config)
		// Connection string
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			options.DbHost, options.DbPort, options.DbUser, options.DbPass, options.DbName)

		// Connect to PostgreSQL
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err.Error())
		}

		userController := controllers.UserController{DB: db}
		// Register user endpoints
		huma.Register(api, huma.Operation{
			OperationID: "register-user",
			Method:      http.MethodPost,
			Path:        "/api/register",
			Summary:     "Register a new user",
			Description: "Registers a new user by accepting email and password.",
			Tags:        []string{"Users"},
		}, userController.CreateUser)

		huma.Register(api, huma.Operation{
			OperationID: "update-user",
			Method:      http.MethodPut,
			Path:        "/api/users/{id}",
			Summary:     "Update a user's account details",
			Description: "Updates the email and/or password for the user specified by the `id`.",
			Tags:        []string{"Users"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, userController.UpdateUser)

		// Register agenda source endpoints
		huma.Register(api, huma.Operation{
			OperationID: "get-agenda-sources",
			Method:      http.MethodGet,
			Path:        "/api/agenda-sources",
			Summary:     "Get a list of agenda sources",
			Description: "Retrieves a list of all agenda sources. Supports ordering by `updatedAt` and pagination.",
			Tags:        []string{"Agenda Sources"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *GetAgendaSourcesInput) (*GetAgendaSourcesOutput, error) {
			// This is a mock implementation
			resp := &GetAgendaSourcesOutput{}
			resp.Body.Data = []AgendaSource{
				{
					ID:        uuid.New().String(),
					URL:       "https://example.com/calendar",
					Type:      "proton",
					UserID:    uuid.New().String(),
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now(),
				},
			}
			resp.Body.Pagination = Pagination{
				Page:       input.Page,
				PageSize:   input.PageSize,
				TotalItems: 1,
				TotalPages: 1,
			}
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "create-agenda-source",
			Method:      http.MethodPost,
			Path:        "/api/agenda-sources",
			Summary:     "Create a new agenda source",
			Description: "Creates a new agenda source with a URL and type.",
			Tags:        []string{"Agenda Sources"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *CreateAgendaSourceInput) (*CreateAgendaSourceOutput, error) {
			// This is a mock implementation
			resp := &CreateAgendaSourceOutput{}
			resp.Body.ID = uuid.New().String()
			resp.Body.URL = input.Body.URL
			resp.Body.Type = input.Body.Type
			resp.Body.UserID = uuid.New().String()
			resp.Body.CreatedAt = time.Now()
			resp.Body.UpdatedAt = time.Now()
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "get-agenda-source",
			Method:      http.MethodGet,
			Path:        "/api/agenda-sources/{id}",
			Summary:     "Get an agenda source by ID",
			Description: "Retrieves an agenda source by its unique identifier.",
			Tags:        []string{"Agenda Sources"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *GetAgendaSourceInput) (*GetAgendaSourceOutput, error) {
			// This is a mock implementation
			resp := &GetAgendaSourceOutput{}
			resp.Body.ID = input.ID
			resp.Body.URL = "https://example.com/calendar"
			resp.Body.Type = "proton"
			resp.Body.UserID = uuid.New().String()
			resp.Body.CreatedAt = time.Now().Add(-24 * time.Hour)
			resp.Body.UpdatedAt = time.Now()
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "update-agenda-source",
			Method:      http.MethodPut,
			Path:        "/api/agenda-sources/{id}",
			Summary:     "Update an agenda source by ID",
			Description: "Updates the URL and/or type of an agenda source by its unique identifier.",
			Tags:        []string{"Agenda Sources"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *UpdateAgendaSourceInput) (*UpdateAgendaSourceOutput, error) {
			// This is a mock implementation
			resp := &UpdateAgendaSourceOutput{}
			resp.Body.ID = input.ID
			resp.Body.URL = input.Body.URL
			resp.Body.Type = input.Body.Type
			resp.Body.UserID = uuid.New().String()
			resp.Body.CreatedAt = time.Now().Add(-24 * time.Hour)
			resp.Body.UpdatedAt = time.Now()
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "delete-agenda-source",
			Method:      http.MethodDelete,
			Path:        "/api/agenda-sources/{id}",
			Summary:     "Delete an agenda source by ID",
			Description: "Deletes an agenda source by its unique identifier.",
			Tags:        []string{"Agenda Sources"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *DeleteAgendaSourceInput) (*struct{}, error) {
			// This is a mock implementation - just return 204 No Content
			return nil, nil
		})

		// Register agenda item endpoints
		huma.Register(api, huma.Operation{
			OperationID: "create-agenda-items",
			Method:      http.MethodPost,
			Path:        "/api/agenda-items",
			Summary:     "Create or update multiple agenda items",
			Description: "Accepts multiple AgendaItem objects and creates or updates them.",
			Tags:        []string{"Agenda Items"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *CreateAgendaItemsInput) (*struct{}, error) {
			// This is a mock implementation - just return 200 OK
			return nil, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "get-agenda-items",
			Method:      http.MethodGet,
			Path:        "/api/agenda-items",
			Summary:     "Query agenda items",
			Description: "Retrieves agenda items based on query parameters.",
			Tags:        []string{"Agenda Items"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *GetAgendaItemsInput) (*GetAgendaItemsOutput, error) {
			// This is a mock implementation
			resp := &GetAgendaItemsOutput{}
			resp.Body.Data = []AgendaItem{
				{
					ResourceID:     uuid.New().String(),
					StartTime:      time.Now().Add(1 * time.Hour),
					EndTime:        time.Now().Add(2 * time.Hour),
					Description:    "Sample agenda item",
					AgendaSourceID: uuid.New().String(),
					UserID:         uuid.New().String(),
				},
			}
			resp.Body.Pagination = Pagination{
				Page:       input.Page,
				PageSize:   input.PageSize,
				TotalItems: 1,
				TotalPages: 1,
			}
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "get-agenda-item",
			Method:      http.MethodGet,
			Path:        "/api/agenda-items/{id}",
			Summary:     "Get an agenda item by ID",
			Description: "Retrieves an agenda item by its ResourceID.",
			Tags:        []string{"Agenda Items"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *GetAgendaItemInput) (*GetAgendaItemOutput, error) {
			// This is a mock implementation
			resp := &GetAgendaItemOutput{}
			resp.Body.ResourceID = input.ID
			resp.Body.StartTime = time.Now().Add(1 * time.Hour)
			resp.Body.EndTime = time.Now().Add(2 * time.Hour)
			resp.Body.Description = "Sample agenda item"
			resp.Body.AgendaSourceID = uuid.New().String()
			resp.Body.UserID = uuid.New().String()
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "delete-agenda-item",
			Method:      http.MethodDelete,
			Path:        "/api/agenda-items/{id}",
			Summary:     "Delete one or multiple agenda items",
			Description: "Deletes agenda items by their ResourceIDs.",
			Tags:        []string{"Agenda Items"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *DeleteAgendaItemInput) (*struct{}, error) {
			// This is a mock implementation - just return 204 No Content
			return nil, nil
		})

		// Register agenda invite endpoints
		huma.Register(api, huma.Operation{
			OperationID: "create-agenda-invite",
			Method:      http.MethodPost,
			Path:        "/api/agenda-invites",
			Summary:     "Create a new agenda invite",
			Description: "Creates a new AgendaInvite.",
			Tags:        []string{"Agenda Invites"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *CreateAgendaInviteInput) (*CreateAgendaInviteOutput, error) {
			// This is a mock implementation
			resp := &CreateAgendaInviteOutput{}
			resp.Body = input.Body
			resp.Body.ResourceID = uuid.New().String()
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "get-agenda-invite",
			Method:      http.MethodGet,
			Path:        "/api/agenda-invites/{id}",
			Summary:     "Get an agenda invite by ID",
			Description: "Retrieves an AgendaInvite by its ResourceID.",
			Tags:        []string{"Agenda Invites"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *GetAgendaInviteInput) (*GetAgendaInviteOutput, error) {
			// This is a mock implementation
			resp := &GetAgendaInviteOutput{}
			resp.Body.ResourceID = input.ID
			resp.Body.UserID = uuid.New().String()
			resp.Body.Description = "Sample agenda invite"
			resp.Body.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
			resp.Body.NotBefore = time.Now()
			resp.Body.NotAfter = time.Now().Add(7 * 24 * time.Hour)
			resp.Body.PaddingBefore = "30m"
			resp.Body.PaddingAfter = "30m"
			resp.Body.SlotSizes = []string{"1h", "30m"}
			resp.Body.AgendaSources = []AgendaSource{
				{
					ID:        uuid.New().String(),
					URL:       "https://example.com/calendar",
					Type:      "proton",
					UserID:    uuid.New().String(),
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now(),
				},
			}
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "update-agenda-invite",
			Method:      http.MethodPut,
			Path:        "/api/agenda-invites/{id}",
			Summary:     "Update an agenda invite by ID",
			Description: "Updates an AgendaInvite by its ResourceID.",
			Tags:        []string{"Agenda Invites"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *UpdateAgendaInviteInput) (*UpdateAgendaInviteOutput, error) {
			// This is a mock implementation
			resp := &UpdateAgendaInviteOutput{}
			resp.Body = input.Body
			resp.Body.ResourceID = input.ID
			return resp, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "delete-agenda-invite",
			Method:      http.MethodDelete,
			Path:        "/api/agenda-invites/{id}",
			Summary:     "Delete an agenda invite by ID",
			Description: "Deletes an AgendaInvite by its ResourceID.",
			Tags:        []string{"Agenda Invites"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		}, func(ctx context.Context, input *DeleteAgendaInviteInput) (*struct{}, error) {
			// This is a mock implementation - just return 204 No Content
			return nil, nil
		})

		huma.Register(api, huma.Operation{
			OperationID: "view-agenda-invite",
			Method:      http.MethodGet,
			Path:        "/api/view-agenda-invite/{id}",
			Summary:     "Publicly available view of a user agenda",
			Description: "Retrieves a list of AgendaItemViews for the specified invite ID within the given date range.",
			Tags:        []string{"Agenda Invites"},
		}, func(ctx context.Context, input *ViewAgendaInviteInput) (*ViewAgendaInviteOutput, error) {
			// This is a mock implementation
			resp := &ViewAgendaInviteOutput{}
			resp.Body = []AgendaItemView{
				{
					StartTime:   time.Now().Add(1 * time.Hour),
					EndTime:     time.Now().Add(2 * time.Hour),
					Description: "Sample agenda item",
				},
				{
					StartTime:   time.Now().Add(3 * time.Hour),
					EndTime:     time.Now().Add(4 * time.Hour),
					Description: "Another sample agenda item",
				},
			}
			return resp, nil
		})

		// Tell the CLI how to start the router
		hooks.OnStart(func() {
			fmt.Printf("Server started on port %d\n", options.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
			if err != nil {
				log.Fatalf("Error starting server: %v", err)
			}

		})
	})

	// Run the CLI
	cli.Run()
}
