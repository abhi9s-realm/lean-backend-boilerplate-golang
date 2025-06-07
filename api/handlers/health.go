package handlers

import (
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, nil, "Service is healthy")
}
