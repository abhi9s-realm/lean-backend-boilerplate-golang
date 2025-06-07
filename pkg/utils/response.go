package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}
