package middleware

import (
	"net/http"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimiter 限流器
type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// Limit 限流中间件
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:ip:" + ip

		pipe := rl.client.Pipeline()
		incr := pipe.Incr(c.Request.Context(), key)
		pipe.Expire(c.Request.Context(), key, rl.window)
		_, err := pipe.Exec(c.Request.Context())
		if err != nil {
			logger.ErrorLog("限流检查失败", logger.Error(err))
			c.Next()
			return
		}

		count, _ := incr.Result()
		if count > int64(rl.limit) {
			c.JSON(http.StatusTooManyRequests, model.APIResponse{
				Code:    "RATE_LIMITED",
				Message: "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
