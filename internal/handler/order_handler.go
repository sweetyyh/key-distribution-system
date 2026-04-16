package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"key-distribution-system/internal/middleware"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

type CreateOrderReq struct {
	ProductID  uint64 `json:"product_id"  binding:"required"`
	Quantity   int    `json:"quantity"    binding:"required,min=1"`
	PayChannel string `json:"pay_channel" binding:"required"`
	Email      string `json:"email"       binding:"required,email"`
	Name       string `json:"name"`
}

// CreateOrder POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid request: "+err.Error())
		return
	}

	var userID uint64
	if id, exists := c.Get(middleware.CtxUserID); exists {
		userID, _ = id.(uint64)
	}

	out, err := store.Global.CreateOrder(store.CreateOrderInput{
		UserID:     userID,
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
			response.Error(c, http.StatusBadRequest, 40003, err.Error())
		}
		return
	}

	response.OK(c, out)
}

// GetOrderDetail GET /api/v1/orders/:order_no
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	orderNo := c.Param("order_no")
	userID, _ := c.Get(middleware.CtxUserID)
	uid, _ := userID.(uint64)

	dto, err := store.Global.GetOrderDetail(orderNo, uid)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40403, "order not found")
		return
	}
	response.OK(c, dto)
}

// ListOrders GET /api/v1/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, _ := c.Get(middleware.CtxUserID)
	uid, _ := userID.(uint64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := store.Global.ListOrders(uid, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "query failed")
		return
	}
	response.OK(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetOrderCards GET /api/v1/orders/:order_no/cards
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
