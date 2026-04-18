package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"key-distribution-system/internal/config"
	"key-distribution-system/internal/pkg/epay"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

// PayHandler 支付相关处理器
type PayHandler struct {
	cfg *config.Config
}

// NewPayHandler 创建支付处理器
func NewPayHandler(cfg *config.Config) *PayHandler {
	return &PayHandler{cfg: cfg}
}

// findChannel 按 name 查找支付渠道配置
func (h *PayHandler) findChannel(name string) (*config.ChannelConfig, bool) {
	for i := range h.cfg.Payment.Channels {
		ch := &h.cfg.Payment.Channels[i]
		if ch.Name == name && ch.Enabled {
			return ch, true
		}
	}
	return nil, false
}

// CreatePayReq 发起支付请求体
type CreatePayReq struct {
	OrderNo  string `json:"order_no"  binding:"required"`
	PayType  string `json:"pay_type"  binding:"required"` // alipay / wxpay / qqpay
	SiteName string `json:"site_name"`
}

// CreatePay POST /api/v1/pay/create/:channel
// 返回支付跳转 URL，前端拿到后直接跳转
func (h *PayHandler) CreatePay(c *gin.Context) {
	channel := c.Param("channel")
	ch, ok := h.findChannel(channel)
	if !ok {
		response.Error(c, http.StatusBadRequest, 40030, "payment channel not found or disabled")
		return
	}

	var req CreatePayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40031, "invalid request: "+err.Error())
		return
	}

	// 从数据库取订单信息
	order, err := store.Global.GetOrderByNo(req.OrderNo)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40032, "order not found")
		return
	}
	if order.Status != 0 {
		response.Error(c, http.StatusBadRequest, 40033, "order already paid or expired")
		return
	}

	siteName := req.SiteName
	if siteName == "" {
		siteName = h.cfg.App.Name
	}

	// 用实际 return_url 替换 {order_no} 占位符
	returnURL := strings.ReplaceAll(ch.ReturnURL, "{order_no}", req.OrderNo)

	payURL := epay.BuildPayURL(epay.PayParams{
		Money:      order.TotalAmount,
		Name:       fmt.Sprintf("Order %s", req.OrderNo),
		NotifyURL:  ch.NotifyURL,
		OutTradeNo: req.OrderNo,
		PID:        ch.PID,
		ReturnURL:  returnURL,
		SiteName:   siteName,
		Type:       req.PayType,
		Key:        ch.Key,
		SubmitURL:  ch.APIURL,
	})

	response.OK(c, gin.H{"pay_url": payURL})
}

// Callback GET /api/v1/pay/callback/:channel
// 接收支付平台服务器异步通知（GET 参数），验签后完成订单
func (h *PayHandler) Callback(c *gin.Context) {
	channel := c.Param("channel")
	ch, ok := h.findChannel(channel)
	if !ok {
		c.String(http.StatusBadRequest, "channel not found")
		return
	}

	// 收集所有 query 参数
	params := make(map[string]string)
	for k, vs := range c.Request.URL.Query() {
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}

	if !epay.VerifyCallback(params, ch.Key) {
		c.String(http.StatusBadRequest, "sign error")
		return
	}

	// 易支付只在 trade_status=TRADE_SUCCESS 时才算成功
	if params["trade_status"] != "TRADE_SUCCESS" {
		c.String(http.StatusOK, "fail")
		return
	}

	orderNo := params["out_trade_no"]
	if err := store.Global.FulfillOrder(orderNo); err != nil {
		switch err.Error() {
		case "order expired":
			c.String(http.StatusOK, "fail")
		default:
			c.String(http.StatusOK, "fail")
		}
		return
	}

	// 易支付要求成功时返回字符串 "success"
	c.String(http.StatusOK, "success")
}

// Return GET /api/v1/pay/return/:channel
// 用户支付完成后从支付平台跳回，验签并重定向到订单页
func (h *PayHandler) Return(c *gin.Context) {
	channel := c.Param("channel")
	ch, ok := h.findChannel(channel)
	if !ok {
		response.Error(c, http.StatusBadRequest, 40034, "channel not found")
		return
	}

	params := make(map[string]string)
	for k, vs := range c.Request.URL.Query() {
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}

	if !epay.VerifyCallback(params, ch.Key) {
		response.Error(c, http.StatusBadRequest, 40035, "sign error")
		return
	}

	orderNo := params["out_trade_no"]
	// 重定向到前端订单详情页
	redirectURL := strings.ReplaceAll(ch.ReturnURL, "{order_no}", orderNo)
	c.Redirect(http.StatusFound, redirectURL)
}
