package main

import (
	"bicycle-store/internal/config"
	"bicycle-store/internal/database"
	"bicycle-store/internal/middleware"
	"bicycle-store/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

// @title Sports Bicycle Store API
// @version 1.0
// @description REST API for Sports Bicycle Store Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bicyclestore.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to MongoDB
	if err := database.Connect(cfg.MongoURI, cfg.DBName); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Create Gin router
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(middleware.RecoveryHandler())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Bicycle Store API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
