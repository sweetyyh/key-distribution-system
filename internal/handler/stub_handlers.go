package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
)

func NotImplemented(c *gin.Context) {
	response.Error(c, http.StatusNotImplemented, 10001, "not implemented")
}
