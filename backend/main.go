package main

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=codegen-config.yaml openapi.yaml

import (
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"awesomeProject/controllers"
	api "awesomeProject/openAPIGenerated"
	"awesomeProject/routers"
)

func main() {
	// Define CLI arguments
	host := flag.String("host", "localhost", "Database host")
	port := flag.String("port", "5432", "Database port")
	user := flag.String("user", "postgres", "Database user")
	dbname := flag.String("dbname", "postgres", "Database name")

	// Parse the command line arguments
	flag.Parse()

	// Get password from environment variable
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	// Create a new Gin router
	r := gin.Default()

	// Set up CORS middleware if needed
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", *host, *port, *user, password, *dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database at %s:%s", *host, *port)

	// Initialize controllers
	userController := controllers.UserController{DB: db}

	// Create server with controllers
	server := routers.NewServer(&userController)

	// Register API handlers
	api.RegisterHandlers(r, server)

	// Start the server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
