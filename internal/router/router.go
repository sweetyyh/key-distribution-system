package router

import (
	"github.com/gin-gonic/gin"
	"key-distribution-system/internal/handler"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	health := handler.NewHealthHandler()
	r.GET("/ping", health.Ping)

	apiV1 := r.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", handler.NotImplemented)
			auth.POST("/login", handler.NotImplemented)
		}

		apiV1.GET("/products", handler.NotImplemented)
		apiV1.GET("/products/:id", handler.NotImplemented)
		apiV1.POST("/orders", handler.NotImplemented)
		apiV1.GET("/orders", handler.NotImplemented)
		apiV1.GET("/orders/:order_no", handler.NotImplemented)
		apiV1.GET("/orders/:order_no/cards", handler.NotImplemented)

		pay := apiV1.Group("/pay")
		{
			pay.POST("/callback/:channel", handler.NotImplemented)
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
