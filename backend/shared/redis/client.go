package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/redis/go-redis/v9"
)

// Client 包装 go-redis，供会话 / 限流 / 分布式锁试点使用（T-603）。
type Client struct {
	rdb *redis.Client
}

// New 根据配置创建 Redis 客户端；addr 为空时返回 (nil, nil) 表示未启用。
func New(cfg config.RedisConfig) (*Client, error) {
	if cfg.Host == "" {
		return nil, nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

func (c *Client) Underlying() *redis.Client {
	if c == nil {
		return nil
	}
	return c.rdb
}

func (c *Client) Close() error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}

// Allow fixed-window 限流：key 在 window 内不超过 limit 次则返回 true。
func (c *Client) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	if c == nil || c.rdb == nil {
		return true, nil
	}
	pipe := c.rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return incr.Val() <= int64(limit), nil
}
