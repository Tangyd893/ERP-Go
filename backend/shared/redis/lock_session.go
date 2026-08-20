package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Lock 基于 SET NX EX 的简单分布式锁（T-603）。
type Lock struct {
	client *Client
	key    string
	token  string
	ttl    time.Duration
}

// AcquireLock 尝试获取锁；成功返回 Lock，失败返回 (nil, nil) 表示未抢到。
func (c *Client) AcquireLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	if c == nil || c.rdb == nil {
		return nil, fmt.Errorf("redis not configured")
	}
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	token := uuid.New().String()
	ok, err := c.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &Lock{client: c, key: key, token: token, ttl: ttl}, nil
}

// Unlock 仅在持有 token 时释放（Lua 防误删）。
func (l *Lock) Unlock(ctx context.Context) error {
	if l == nil || l.client == nil || l.client.rdb == nil {
		return nil
	}
	const script = `
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end`
	return l.client.rdb.Eval(ctx, script, []string{l.key}, l.token).Err()
}

// SessionStore 简单会话 KV（登录态试点，非完整 IAM 会话表）。
type SessionStore struct {
	client *Client
	prefix string
}

func NewSessionStore(c *Client) *SessionStore {
	return &SessionStore{client: c, prefix: "sess:"}
}

func (s *SessionStore) Put(ctx context.Context, sessionID, userJSON string, ttl time.Duration) error {
	if s == nil || s.client == nil || s.client.rdb == nil {
		return fmt.Errorf("redis not configured")
	}
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return s.client.rdb.Set(ctx, s.prefix+sessionID, userJSON, ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, sessionID string) (string, error) {
	if s == nil || s.client == nil || s.client.rdb == nil {
		return "", fmt.Errorf("redis not configured")
	}
	val, err := s.client.rdb.Get(ctx, s.prefix+sessionID).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	if s == nil || s.client == nil || s.client.rdb == nil {
		return nil
	}
	return s.client.rdb.Del(ctx, s.prefix+sessionID).Err()
}
