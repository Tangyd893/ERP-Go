//go:build integration

package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func setupOrderTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	// 优先使用 TEST_DATABASE_URL（外部 PG，无需 Docker）
	if connStr := getEnv("TEST_DATABASE_URL"); connStr != "" {
		db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Warn),
		})
		require.NoError(t, err, "连接外部 PostgreSQL 失败")
		if err := autoMigrateOrder(db); err != nil {
			t.Logf("自动建表警告: %v", err)
		}
		return db, func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}

	// 尝试 testcontainers（需 Docker）
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
		t.Skipf("跳过 DB 集成测试（Docker 不可用）: %v\n  提示: 设置 TEST_DATABASE_URL 环境变量可使用外部 PostgreSQL", err)
		return nil, func() {}
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	require.NoError(t, err, "连接数据库失败")

	err = autoMigrateOrder(db)
	require.NoError(t, err, "自动建表失败")

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func autoMigrateOrder(db *gorm.DB) error {
	return db.AutoMigrate(
		&SalesOrderModel{},
		&OrderItemModel{},
		&OrderAddressModel{},
		&OrderStatusLogModel{},
	)
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func newTestOrder(tenantID, orderNo string) *domain.SalesOrder {
	now := time.Now()
	return &domain.SalesOrder{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		StoreID:         "store-001",
		PlatformOrderNo: orderNo,
		OrderType:       domain.OrderTypeNormal,
		OrderSource:     domain.OrderSourcePlatform,
		Status:          domain.OrderPending,
		BuyerName:       "测试买家",
		BuyerEmail:      "test@example.com",
		Currency:        "CNY",
		TotalAmount:     99.90,
		ShippingFee:     10.00,
		TaxAmount:       5.00,
		Items: []*domain.OrderItem{
			{ID: uuid.New().String(), SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Quantity: 2, UnitPrice: 49.95, TotalPrice: 99.90},
		},
		Address: &domain.Address{
			ContactName: "张三", Phone: "13800138000", Email: "zhang@example.com",
			Country: "中国", State: "浙江", City: "杭州", District: "余杭区",
			StreetLine1: "文一西路 969 号", PostalCode: "311121",
		},
		IdempotencyKey: uuid.New().String(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// TestOrderRepositoryCRUD 验证订单仓储 CRUD（真实 PostgreSQL）
func TestOrderRepositoryCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 DB 集成测试 (go test -short)")
	}

	db, cleanup := setupOrderTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(db)
	ctx := context.Background()
	tenantID := "t-dbtest-001"

	order := newTestOrder(tenantID, "DBTEST-001")
	err := repo.Create(ctx, order)
	require.NoError(t, err, "创建订单失败")

	found, err := repo.FindByID(ctx, order.ID)
	require.NoError(t, err, "查询订单失败")
	assert.Equal(t, order.PlatformOrderNo, found.PlatformOrderNo)
	assert.Equal(t, order.BuyerName, found.BuyerName)
	assert.Equal(t, order.TotalAmount, found.TotalAmount)
	assert.Len(t, found.Items, 1)
	assert.Equal(t, "sku-001", found.Items[0].SKUID)
	assert.NotNil(t, found.Address)
	assert.Equal(t, "张三", found.Address.ContactName)

	err = repo.UpdateStatus(ctx, order.ID, "approved")
	require.NoError(t, err, "更新状态失败")
	found, err = repo.FindByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderApproved, found.Status)

	dup, err := repo.FindByIdempotencyKey(ctx, order.IdempotencyKey)
	require.NoError(t, err, "按幂等键查询失败")
	assert.Equal(t, order.ID, dup.ID)

	order2 := newTestOrder(tenantID, "DBTEST-002")
	order2.IdempotencyKey = uuid.New().String()
	require.NoError(t, repo.Create(ctx, order2))

	orders, total, err := repo.List(ctx, tenantID, 0, 10)
	require.NoError(t, err, "列表查询失败")
	assert.Equal(t, int64(2), total)
	assert.Len(t, orders, 2)
	assert.Equal(t, "DBTEST-002", orders[0].PlatformOrderNo)
	assert.Equal(t, "DBTEST-001", orders[1].PlatformOrderNo)
}

func TestOrderRepositoryFindByIDNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 DB 集成测试")
	}
	db, cleanup := setupOrderTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(db)
	_, err := repo.FindByID(context.Background(), "nonexistent-id")
	assert.Error(t, err, "应返回错误")
}

func TestOrderRepositoryIdempotency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 DB 集成测试")
	}
	db, cleanup := setupOrderTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(db)
	ctx := context.Background()

	key := uuid.New().String()
	order1 := newTestOrder("t-dup", "DUP-001")
	order1.IdempotencyKey = key
	require.NoError(t, repo.Create(ctx, order1))

	order2 := newTestOrder("t-dup", "DUP-002")
	order2.IdempotencyKey = key
	err := repo.Create(ctx, order2)
	assert.Error(t, err, "重复幂等键应报错")
}
