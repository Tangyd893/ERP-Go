package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	httpiface "github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/gin-gonic/gin"
)

func setupOrderTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	return engine
}

func TestOrderHandler_FallbackList(t *testing.T) {
	engine := setupOrderTestEngine()

	// fallback 模式：appService 为 nil
	handler := httpiface.NewOrderHandler(nil)
	router := engine.Group("/api/v1")
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 0 {
		t.Errorf("期望 code=0，实际 %v", resp["code"])
	}
}

func TestOrderHandler_FallbackGet(t *testing.T) {
	engine := setupOrderTestEngine()

	handler := httpiface.NewOrderHandler(nil)
	router := engine.Group("/api/v1")
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/any-id", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("期望 404，实际 %d", w.Code)
	}
}

func TestOrderHandler_RealList(t *testing.T) {
	engine := setupOrderTestEngine()

	// 使用真实 MemOutboxStore + mock repo
	store := outbox.NewMemOutboxStore()
	appSvc := app.NewOrderAppService(newMockOrderRepo()).WithOutbox(store)
	handler := httpiface.NewOrderHandler(appSvc)
	router := engine.Group("/api/v1")
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 0 {
		t.Errorf("期望 code=0，实际 %v", resp["code"])
	}
}

func TestOrderHandler_OutboxEndpoints(t *testing.T) {
	engine := setupOrderTestEngine()

	store := outbox.NewMemOutboxStore()
	appSvc := app.NewOrderAppService(newMockOrderRepo()).WithOutbox(store)
	handler := httpiface.NewOrderHandler(appSvc).WithOutboxStore(store)
	router := engine.Group("/api/v1")
	handler.RegisterRoutes(router)

	// 查询失败消息列表（空）
	req := httptest.NewRequest(http.MethodGet, "/api/v1/outbox/failed", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d", w.Code)
	}

	// 重试（空队列）
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/outbox/retry", strings.NewReader(`{"id":1}`))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d", w2.Code)
	}
}

// mockOrderRepo 简单内存实现
type mockOrderRepo struct {
	orders map[string]*domain.SalesOrder
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[string]*domain.SalesOrder)}
}

func (m *mockOrderRepo) Create(_ context.Context, order *domain.SalesOrder) error {
	m.orders[order.ID] = order
	return nil
}
func (m *mockOrderRepo) FindByID(_ context.Context, id string) (*domain.SalesOrder, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, nil
}
func (m *mockOrderRepo) UpdateStatus(_ context.Context, id, status string) error {
	if o, ok := m.orders[id]; ok {
		o.Status = domain.OrderStatus(status)
	}
	return nil
}
func (m *mockOrderRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.SalesOrder, int64, error) {
	var list []*domain.SalesOrder
	for _, o := range m.orders {
		list = append(list, o)
	}
	return list, int64(len(list)), nil
}
