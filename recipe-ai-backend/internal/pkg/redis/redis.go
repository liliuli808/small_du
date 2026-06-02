package redis

import (
	"context"
	"time"

	"recipe-ai-backend/internal/pkg/config"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

// Client Redis客户端
type Client struct {
	client *redis.Client
}

// NewClient 创建Redis客户端
func NewClient(cfg config.RedisConfig) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{client: client}
}

// Ping 测试连接
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.client.Close()
}

// Set 设置缓存
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get 获取缓存
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete 删除缓存
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// RateLimitCheck 限流检查
func (c *Client) RateLimitCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	pipe := c.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := cmds[0].(*redis.IntCmd).Result()
	if err != nil {
		return false, err
	}

	return count <= int64(limit), nil
}

// AsynqRedisOpt Asynq Redis选项
func (c *Client) AsynqRedisOpt() asynq.RedisConnOpt {
	return asynq.RedisClientOpt{
		Addr:     c.client.Options().Addr,
		Password: c.client.Options().Password,
		DB:       c.client.Options().DB,
	}
}

// RawClient 返回原始redis客户端
func (c *Client) RawClient() *redis.Client {
	return c.client
}
