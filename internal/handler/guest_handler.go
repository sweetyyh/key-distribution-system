package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// GuestHandler 游客（无需登录）处理器
type GuestHandler struct{}

func NewGuestHandler() *GuestHandler {
	return &GuestHandler{}
}

type GuestCreateOrderReq struct {
	ProductID  uint64 `json:"product_id"  binding:"required"`
	Quantity   int    `json:"quantity"    binding:"required,min=1"`
	PayChannel string `json:"pay_channel" binding:"required"`
	Email      string `json:"email"       binding:"required,email"`
}

// CreateOrder POST /api/v1/guest/orders — 游客创建订单（无需 JWT）
func (h *GuestHandler) CreateOrder(c *gin.Context) {
	var req GuestCreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid request: "+err.Error())
		return
	}

	out, err := store.Global.CreateOrder(store.CreateOrderInput{
		UserID:     0, // 游客，store 会兜底为 1
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		PayChannel: req.PayChannel,
		Email:      req.Email,
	})
	if err != nil {
		switch err.Error() {
		case "product not found":
			response.Error(c, http.StatusNotFound, 40402, "product not found")
		default:
			response.Error(c, http.StatusBadRequest, 40003, err.Error())
		}
		return
	}

	response.Created(c, out)
}

// GetOrderCards GET /api/v1/guest/orders/:order_no?email=xxx — 游客查询订单卡密
func (h *GuestHandler) GetOrderCards(c *gin.Context) {
	orderNo := c.Param("order_no")
	email := c.Query("email")

	dto, err := store.Global.GuestGetOrderCards(orderNo, email)
	if err != nil {
		switch err.Error() {
		case "order not found":
			response.Error(c, http.StatusNotFound, 40403, "order not found")
		case "order not paid yet":
			response.Error(c, http.StatusPaymentRequired, 40005, "order not paid yet")
		case "email mismatch":
			response.Error(c, http.StatusForbidden, 40301, "email mismatch")
		default:
			response.Error(c, http.StatusBadRequest, 40006, err.Error())
		}
		return
	}

	response.OK(c, dto)
}
