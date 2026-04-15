package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/handler"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 托管前端静态文件
	r.Static("/portal", "./stitch_virtual_card_wholesale_portal")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/portal/_1/code.html")
	})

	// CORS 中间件（内联实现，无需额外依赖）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	health := handler.NewHealthHandler()
	r.GET("/ping", health.Ping)

	productH := handler.NewProductHandler()
	orderH := handler.NewOrderHandler()
	payH := handler.NewPayHandler()

	apiV1 := r.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", handler.NotImplemented)
			auth.POST("/login", handler.NotImplemented)
		}

		apiV1.GET("/products", productH.List)
		apiV1.GET("/products/:id", productH.Get)
		apiV1.POST("/orders", orderH.CreateOrder)
		apiV1.GET("/orders", handler.NotImplemented)
		apiV1.GET("/orders/:order_no", handler.NotImplemented)
		apiV1.GET("/orders/:order_no/cards", orderH.GetOrderCards)

		pay := apiV1.Group("/pay")
		{
			pay.POST("/callback/:channel", payH.MockCallback)
			pay.GET("/return/:channel", handler.NotImplemented)
		}
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/products", handler.NotImplemented)
		admin.PUT("/products/:id", handler.NotImplemented)
		admin.POST("/products/:id/cards", handler.NotImplemented)
		admin.GET("/orders", handler.NotImplemented)
		admin.PUT("/orders/:order_no/refund", handler.NotImplemented)
		admin.GET("/orders/failed", handler.NotImplemented)
		admin.GET("/stats", handler.NotImplemented)
	}

	return r
}
