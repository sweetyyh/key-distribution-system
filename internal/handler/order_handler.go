package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// OrderHandler 订单相关处理器
type OrderHandler struct{}

// NewOrderHandler 创建订单处理器
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// CreateOrderReq 创建订单请求体
type CreateOrderReq struct {
	ProductID  uint64 `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	PayChannel string `json:"pay_channel" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name"`
}

// CreateOrder 创建订单 POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid request: "+err.Error())
		return
	}

	out, err := store.Global.CreateOrder(store.CreateOrderInput{
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		PayChannel: req.PayChannel,
		Email:      req.Email,
		Name:       req.Name,
	})
	if err != nil {
		switch err.Error() {
		case "product not found":
			response.Error(c, http.StatusNotFound, 40402, "product not found")
		default:
			// insufficient stock 或其他
			response.Error(c, http.StatusBadRequest, 40003, err.Error())
		}
		return
	}

	response.OK(c, out)
}

// GetOrderCards 获取订单卡密 GET /api/v1/orders/:order_no/cards
func (h *OrderHandler) GetOrderCards(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.Error(c, http.StatusBadRequest, 40004, "order_no is required")
		return
	}

	dto, err := store.Global.GetOrderCards(orderNo)
	if err != nil {
		switch err.Error() {
		case "order not paid yet":
			response.Error(c, http.StatusBadRequest, 40005, "order not paid yet")
		default:
			response.Error(c, http.StatusNotFound, 40403, err.Error())
		}
		return
	}

	response.OK(c, dto)
}
