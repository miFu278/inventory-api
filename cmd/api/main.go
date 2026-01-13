package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"inventory-api/config"
	"inventory-api/database"
	"inventory-api/handler"
	"inventory-api/middleware"
	"inventory-api/models"
	"inventory-api/repo"
	"inventory-api/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.ConnectPostgres(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Auto migrate models
	db := database.GetDB()
	if err := db.AutoMigrate(&models.Product{}, &models.Transaction{}, &models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	// Initialize repositories
	inventoryRepo := repo.NewInventoryRepository(db)
	userRepo := repo.NewUserRepository(db)

	// Initialize services
	inventoryService := services.NewInventoryService(inventoryRepo)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)

	// Initialize handlers
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	userHandler := handler.NewUserHandler(userService)

	// Setup Gin router
	router := gin.Default()

	// Setup Huma API
	humaConfig := huma.DefaultConfig("Inventory API", "1.0.0")
	humaConfig.Info.Description = "REST API for inventory management system with JWT authentication"

	// Add security scheme for JWT
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter your JWT token",
		},
	}

	api := humagin.New(router, humaConfig)

	// Add JWT middleware to protected routes
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		path := ctx.URL().Path

		// Skip auth for public endpoints
		if strings.HasPrefix(path, "/auth/") ||
			strings.HasPrefix(path, "/docs") ||
			strings.HasPrefix(path, "/schemas") ||
			strings.HasPrefix(path, "/openapi") {
			next(ctx)
			return
		}

		// Apply auth middleware for protected routes
		if strings.HasPrefix(path, "/users") ||
			strings.HasPrefix(path, "/products") ||
			strings.HasPrefix(path, "/transactions") {

			// Allow public read access to products list and details
			if (path == "/products" || strings.HasPrefix(path, "/products/")) &&
				ctx.Method() == http.MethodGet &&
				!strings.Contains(path, "/transactions") {
				next(ctx)
				return
			}

			middleware.HumaAuthMiddleware(api, cfg.JWTSecret)(ctx, next)
			return
		}

		next(ctx)
	})

	// Register routes
	inventoryHandler.RegisterRoutes(api)
	userHandler.RegisterRoutes(api)

	// Get server port
	port := cfg.ServerPort

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Printf("API documentation available at http://localhost:%s/docs", port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
