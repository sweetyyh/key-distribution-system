package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// AdminHandler 后台管理处理器
type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetStats GET /api/admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats := store.Global.GetAdminStats()
	response.OK(c, stats)
}

// ListOrders GET /api/admin/orders
func (h *AdminHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	statusStr := c.DefaultQuery("status", "-1")
	status, _ := strconv.Atoi(statusStr)

	items, total, err := store.Global.AdminListOrders(page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "query failed")
		return
	}
	response.OK(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListUsers GET /api/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := store.Global.AdminListUsers(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50004, "query failed")
		return
	}
	response.OK(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
