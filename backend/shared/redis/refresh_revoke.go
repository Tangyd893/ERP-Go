package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const refreshRevokePrefix = "iam:refresh:revoked:"

// RefreshRevokeStore stores revoked refresh-token JTIs with TTL until expiry.
type RefreshRevokeStore struct {
	client *Client
}

func NewRefreshRevokeStore(c *Client) *RefreshRevokeStore {
	if c == nil {
		return nil
	}
	return &RefreshRevokeStore{client: c}
}

func (s *RefreshRevokeStore) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	if s == nil || s.client == nil || s.client.rdb == nil {
		return fmt.Errorf("redis not configured")
	}
	if jti == "" {
		return fmt.Errorf("empty jti")
	}
	if ttl <= 0 {
		ttl = time.Hour
	}
	return s.client.rdb.Set(ctx, refreshRevokePrefix+jti, "1", ttl).Err()
}

func (s *RefreshRevokeStore) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if s == nil || s.client == nil || s.client.rdb == nil {
		return false, fmt.Errorf("redis not configured")
	}
	if jti == "" {
		return false, nil
	}
	n, err := s.client.rdb.Exists(ctx, refreshRevokePrefix+jti).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return n > 0, nil
}
