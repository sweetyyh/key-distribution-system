package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"key-distribution-system/internal/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Ping(c *gin.Context) {
	response.OK(c, gin.H{
		"time": time.Now().UTC().Format(time.RFC3339),
	})
}
