package testutils

import (
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/config"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/services" // New import
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/database"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/persistence" // New import
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getTestConfig returns a test configuration
func getTestConfig() *config.Config {
	return &config.Config{
		Environment: "test",
		Port:        "8081", // Ensure this port does not conflict if tests run in parallel or with dev server
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "postgres", // Ensure these are correct for your test environment
		DBPass:      "postgres",
		DBName:      "test_db", // Ensure this DB exists or can be created by the user
		LogLevel:    "debug",
	}
}

// SetupTestRouter returns a configured Gin router, UserService, and optional database connection for testing
func SetupTestRouter(needsDB bool) (*gin.Engine, services.UserService, *gorm.DB) {
	cfg := getTestConfig()
	logger.NewLogger(cfg.LogLevel) // Initialize logger

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	var db *gorm.DB
	var userService services.UserService

	if needsDB {
		var err error
		db, err = database.NewPostgresDB(cfg)
		if err != nil {
			panic("Failed to connect to test database: " + err.Error())
		}

		// Auto-migrate the database schema for tests
		// Consider if this should be cleared before each test suite or run
		if err := db.AutoMigrate(&models.User{}); err != nil {
			panic("Failed to migrate test database: " + err.Error())
		}

		// Initialize Repository and Service
		userRepository := persistence.NewGormUserRepository(db)
		userService = services.NewUserService(userRepository)
	}

	return r, userService, db
}

// CleanupDatabase cleans up test data from specified tables.
// It is crucial to ensure tests are independent.
func CleanupDatabase(db *gorm.DB) {
	// Add other tables here if necessary
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
        // If TRUNCATE fails (e.g. due to permissions or table locks), fallback or log error
        // Fallback to DELETE for simplicity if TRUNCATE causes issues in some environments
        db.Exec("DELETE FROM users")
    }
}
