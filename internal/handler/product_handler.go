package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// ProductHandler 产品相关处理器
type ProductHandler struct{}

// NewProductHandler 创建产品处理器
func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// List 获取产品列表 GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
	list := store.Global.ListProducts()
	response.OK(c, gin.H{
		"list":  list,
		"total": len(list),
	})
}

// Get 获取产品详情 GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid product id")
		return
	}

	p, ok := store.Global.GetProduct(id)
	if !ok {
		response.Error(c, http.StatusNotFound, 40401, "product not found")
		return
	}

	response.OK(c, p)
}
