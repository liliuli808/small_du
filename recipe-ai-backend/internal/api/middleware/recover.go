package middleware

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Recover 恢复中间件
func Recover() gin.HandlerFunc {
	return gin.RecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		logger.ErrorLog("panic recovered", logger.String("path", c.Request.URL.Path))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    "INTERNAL_ERROR",
			Message: "服务异常，请稍后重试",
		})
	})
}
