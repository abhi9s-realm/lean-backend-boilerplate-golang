package middleware

import (
	"time"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Infow("HTTP Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
