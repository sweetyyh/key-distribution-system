package response

import "github.com/gin-gonic/gin"

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Body{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, statusCode int, code int, msg string) {
	c.JSON(statusCode, Body{Code: code, Message: msg})
}
