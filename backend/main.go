package main

import (
	"awesomeProject/controllers"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
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

		// Create user controller
		userController := &controllers.UserController{DB: db}

		// Register all routes
		addRoutes(api, userController)

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
