package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"key-distribution-system/internal/pkg/jwtutil"
	"key-distribution-system/internal/pkg/response"
)

const (
	CtxUserID = "uid"
	CtxRole   = "role"
)

// BuyerAuth 买家 JWT 鉴权中间件
func BuyerAuth(secret string) gin.HandlerFunc {
	return jwtAuth(secret, "")
}

// AdminAuth admin JWT 鉴权中间件
func AdminAuth(secret string) gin.HandlerFunc {
	return jwtAuth(secret, "admin")
}

func jwtAuth(secret, requireRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, 40100, "missing token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, err := jwtutil.Parse(tokenStr, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40102, "invalid token")
			c.Abort()
			return
		}

		if requireRole != "" && claims.Role != requireRole {
			response.Error(c, http.StatusForbidden, 40301, "forbidden")
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
