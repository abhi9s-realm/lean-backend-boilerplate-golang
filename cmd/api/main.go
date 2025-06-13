package main

import (
	"log"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/routes"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/config"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/services"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/database"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/persistence" // New import
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	l := logger.NewLogger(cfg.LogLevel)
	defer l.Sync() // Ensure logs are flushed

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		// Using l.Fatal with a simple error message as per existing style
		l.Fatal("Failed to connect to database: " + err.Error())
	}

	// Initialize Repositories
	userRepository := persistence.NewGormUserRepository(db)

	// Initialize Services
	userService := services.NewUserService(userRepository)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes - pass userService
	routes.Setup(r, userService, l) // Pass userService instead of db

	// Start server
	l.Info("Starting server on port " + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		// Using l.Fatal with a simple error message as per existing style
		l.Fatal("Failed to start server: " + err.Error())
	}
}
