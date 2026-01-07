package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"inventory-api/internal/database"
	"inventory-api/internal/handler"
	"inventory-api/internal/models"
	"inventory-api/internal/repo"
	"inventory-api/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	if err := database.ConnectPostgres(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Auto migrate models
	db := database.GetDB()
	if err := db.AutoMigrate(&models.Product{}, &models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	// Initialize layers
	inventoryRepo := repo.NewInventoryRepository(db)
	inventoryService := services.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// Setup Gin router
	router := gin.Default()

	// Setup Huma API
	config := huma.DefaultConfig("Inventory API", "1.0.0")
	config.Info.Description = "REST API for inventory management system"
	api := humagin.New(router, config)

	// Register routes
	inventoryHandler.RegisterRoutes(api)

	// Get server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Printf("API documentation available at http://localhost:%s/docs", port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
