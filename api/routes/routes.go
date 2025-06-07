package routes

import (
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/handlers"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/middleware"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, log *logger.Logger) {
	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(log))

	// Health check
	r.GET("/api/health", handlers.HealthCheck)

	// User routes
	userHandler := handlers.NewUserHandler(db)
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}
}
