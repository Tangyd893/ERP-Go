// Package testdb provides shared PostgreSQL helpers for integration tests.
package testdb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// OpenPostgres returns a GORM DB for integration tests.
// Prefers TEST_DATABASE_URL; otherwise starts postgres:16-alpine via testcontainers.
// Skips the test when Docker is unavailable and no URL is set.
func OpenPostgres(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	if connStr := os.Getenv("TEST_DATABASE_URL"); connStr != "" {
		db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Warn),
		})
		require.NoError(t, err, "连接外部 PostgreSQL 失败")
		return db, func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				_ = sqlDB.Close()
			}
		}
	}

	ctx := context.Background()
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("erp_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Skipf("跳过 DB 集成测试（Docker 不可用）: %v\n  提示: 设置 TEST_DATABASE_URL 可使用外部 PostgreSQL", err)
		return nil, func() {}
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	require.NoError(t, err, "连接数据库失败")

	return db, func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		_ = pgContainer.Terminate(ctx)
	}
}
