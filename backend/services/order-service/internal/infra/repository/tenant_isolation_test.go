//go:build integration

package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/app"
	orderhttp "github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/testing/testdb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantIsolation_ListAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 DB 集成测试")
	}

	db, cleanup := testdb.OpenPostgres(t)
	defer cleanup()
	require.NoError(t, autoMigrateOrder(db))

	repo := NewOrderRepository(db)
	ctx := context.Background()

	a := newTestOrder("tenant-a", "ISO-A-001")
	b := newTestOrder("tenant-b", "ISO-B-001")
	require.NoError(t, repo.Create(ctx, a))
	require.NoError(t, repo.Create(ctx, b))

	listA, totalA, err := repo.List(ctx, "tenant-a", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), totalA)
	assert.Equal(t, a.ID, listA[0].ID)

	_, err = repo.FindByID(ctx, "tenant-b", a.ID)
	assert.Error(t, err, "跨租户 FindByID 应失败")

	_, err = repo.FindByID(ctx, "", a.ID)
	assert.Error(t, err, "空租户 FindByID 应 fail-closed")

	err = repo.UpdateStatus(ctx, "", a.ID, "approved")
	assert.Error(t, err, "空租户 UpdateStatus 应 fail-closed")

	svc := app.NewOrderAppService(repo)
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set("tenant_id", c.GetHeader("X-Tenant-ID"))
		c.Next()
	})
	api := engine.Group("/api/v1/order")
	orderhttp.NewOrderHandler(svc).RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/order/orders/"+a.ID, nil)
	req.Header.Set("X-Tenant-ID", "tenant-b")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "跨租户 GetOrder 应 404, body=%s", w.Body.String())
}
