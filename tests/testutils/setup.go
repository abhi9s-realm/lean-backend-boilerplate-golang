package testutils

import (
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/config"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/database"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getTestConfig returns a test configuration
func getTestConfig() *config.Config {
	return &config.Config{
		Environment: "test",
		Port:        "8081",
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "postgres",
		DBPass:      "postgres",
		DBName:      "test_db",
		LogLevel:    "debug",
	}
}

// SetupTestRouter returns a configured Gin router and optional database connection for testing
func SetupTestRouter(needsDB bool) (*gin.Engine, *gorm.DB) {
	// Use test config
	cfg := getTestConfig()

	// Initialize logger
	logger.NewLogger(cfg.LogLevel)

	// Set up router in test mode
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Initialize database only if needed
	var db *gorm.DB
	if needsDB {
		var err error
		db, err = database.NewPostgresDB(cfg)
		if err != nil {
			panic(err)
		}

		// Auto-migrate the database schema
		if err := db.AutoMigrate(&models.User{}); err != nil {
			panic(err)
		}
	}

	return r, db
}

// CleanupDatabase cleans up test data after each test
func CleanupDatabase(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}
