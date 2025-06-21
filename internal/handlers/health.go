package handlers

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		TimeStamp: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}
