package middleware

import (
	"github.com/gin-gonic/gin"
)

// UserOpenID 从请求头中获取用户openid
func UserOpenID() gin.HandlerFunc {
	return func(c *gin.Context) {
		openid := c.GetHeader("X-User-OpenID")
		if openid != "" {
			c.Set("user_openid", openid)
		}
		c.Next()
	}
}
