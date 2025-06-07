package main

import (
	"log"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/routes"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/config"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/database"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
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
	defer l.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		l.Fatal("Failed to connect to database", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.Setup(r, db, l)

	// Start server
	l.Info("Starting server on port " + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		l.Fatal("Failed to start server", err)
	}
}
