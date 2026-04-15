package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// PayHandler 支付相关处理器
type PayHandler struct{}

// NewPayHandler 创建支付处理器
func NewPayHandler() *PayHandler {
	return &PayHandler{}
}

// MockCallbackReq 模拟支付回调请求体
type MockCallbackReq struct {
	OrderNo string `json:"order_no" binding:"required"`
}

// MockCallback 模拟支付完成回调 POST /api/v1/pay/callback/:channel
// 在 mock 场景下由前端直接调用，立即完成支付流程
func (h *PayHandler) MockCallback(c *gin.Context) {
	var req MockCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40010, "invalid request: "+err.Error())
		return
	}

	if err := store.Global.FulfillOrder(req.OrderNo); err != nil {
		switch err.Error() {
		case "order expired":
			response.Error(c, http.StatusBadRequest, 40011, "order expired")
		default:
			response.Error(c, http.StatusNotFound, 40404, err.Error())
		}
		return
	}

	response.OK(c, gin.H{"status": "paid"})
}
