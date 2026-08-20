package redisx

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) (*Client, func()) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := &Client{rdb: rdb}
	return c, func() {
		_ = rdb.Close()
		mr.Close()
	}
}

func TestAcquireLock_AndUnlock(t *testing.T) {
	c, cleanup := testClient(t)
	defer cleanup()
	ctx := context.Background()

	lock, err := c.AcquireLock(ctx, "lock:test", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, lock)

	again, err := c.AcquireLock(ctx, "lock:test", time.Minute)
	require.NoError(t, err)
	require.Nil(t, again, "同一 key 不应重复获锁")

	require.NoError(t, lock.Unlock(ctx))
	again, err = c.AcquireLock(ctx, "lock:test", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, again)
	_ = again.Unlock(ctx)
}

func TestSessionStore_PutGetDelete(t *testing.T) {
	c, cleanup := testClient(t)
	defer cleanup()
	ctx := context.Background()
	ss := NewSessionStore(c)

	require.NoError(t, ss.Put(ctx, "s1", `{"user":"a"}`, time.Minute))
	val, err := ss.Get(ctx, "s1")
	require.NoError(t, err)
	require.Equal(t, `{"user":"a"}`, val)
	require.NoError(t, ss.Delete(ctx, "s1"))
	val, err = ss.Get(ctx, "s1")
	require.NoError(t, err)
	require.Empty(t, val)
}
