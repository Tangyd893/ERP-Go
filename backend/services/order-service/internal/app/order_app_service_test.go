package app

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	sharedEvents "github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

// mockOrderRepo 订单仓储的内存模拟实现
type mockOrderRepo struct {
	orders map[string]*domain.SalesOrder
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[string]*domain.SalesOrder)}
}

func (m *mockOrderRepo) Create(ctx context.Context, order *domain.SalesOrder) error {
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*domain.SalesOrder, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id, status string) error {
	if o, ok := m.orders[id]; ok {
		o.Status = domain.OrderStatus(status)
	}
	return nil
}

func (m *mockOrderRepo) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SalesOrder, int64, error) {
	return nil, 0, nil
}

// assertOutboxMessages 验证 outbox 中的消息数量和消息属性
func assertPendingMessages(t *testing.T, store outbox.OutboxStore, wantCount int) []*outbox.OutboxMessage {
	t.Helper()
	msgs, err := store.FetchPending(context.Background(), 100)
	if err != nil {
		t.Fatalf("读取 outbox 失败: %v", err)
	}
	if len(msgs) != wantCount {
		t.Fatalf("期望 %d 条 pending 消息，实际 %d 条", wantCount, len(msgs))
	}
	return msgs
}

// assertEventPayload 解析 payload 并断言关键字段
func assertEventPayload(t *testing.T, msg *outbox.OutboxMessage, wantEventType string, wantAggregateID string) map[string]interface{} {
	t.Helper()
	if msg.EventType != wantEventType {
		t.Errorf("期望事件类型 %s，实际 %s", wantEventType, msg.EventType)
	}
	if msg.AggregateID != wantAggregateID {
		t.Errorf("期望聚合ID %s，实际 %s", wantAggregateID, msg.AggregateID)
	}
	if msg.AggregateType != "SalesOrder" {
		t.Errorf("期望聚合类型 SalesOrder，实际 %s", msg.AggregateType)
	}
	if msg.TenantID != "" {
		t.Logf("事件消息租户ID: %s", msg.TenantID)
	}
	if msg.Status != outbox.StatusPending {
		t.Errorf("期望消息状态 pending，实际 %s", msg.Status)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		t.Fatalf("解析 payload 失败: %v", err)
	}
	if et, ok := payload["event_type"].(string); !ok || et != wantEventType {
		t.Errorf("payload 中 event_type 期望 %s，实际 %v", wantEventType, payload["event_type"])
	}
	return payload
}

// TestCreateOrderEmitsEvent 验证创建订单时正确发布 order.imported 事件
func TestCreateOrderEmitsEvent(t *testing.T) {
	repo := newMockOrderRepo()
	outboxStore := outbox.NewMemOutboxStore()
	svc := NewOrderAppService(repo).WithOutbox(outboxStore)

	order := &domain.SalesOrder{
		ID:             "order-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		PlatformOrderNo: "PO-2024-001",
		Status:          domain.OrderPending,
		IdempotencyKey:  "ik-001",
	}
	if err := svc.CreateOrder(context.Background(), order); err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}

	msgs := assertPendingMessages(t, outboxStore, 1)
	assertEventPayload(t, msgs[0], sharedEvents.EventOrderImported, "order-001")
	if msgs[0].TenantID != "tenant-001" {
		t.Errorf("期望租户ID tenant-001，实际 %s", msgs[0].TenantID)
	}
}

// TestApproveOrderEmitsEvent 验证审核订单时正确发布 order.approved 事件
func TestApproveOrderEmitsEvent(t *testing.T) {
	repo := newMockOrderRepo()
	outboxStore := outbox.NewMemOutboxStore()
	svc := NewOrderAppService(repo).WithOutbox(outboxStore)

	order := &domain.SalesOrder{
		ID:             "order-002",
		TenantID:       "tenant-002",
		StoreID:        "store-001",
		PlatformOrderNo: "PO-2024-002",
		Status:          domain.OrderPending,
		IdempotencyKey:  "ik-002",
	}
	repo.orders[order.ID] = order

	if err := svc.ApproveOrder(context.Background(), "order-002", "admin"); err != nil {
		t.Fatalf("审核订单失败: %v", err)
	}

	msgs := assertPendingMessages(t, outboxStore, 1)
	assertEventPayload(t, msgs[0], sharedEvents.EventOrderApproved, "order-002")

	if msgs[0].TenantID != "tenant-002" {
		t.Errorf("期望租户ID tenant-002，实际 %s", msgs[0].TenantID)
	}

	if updated, ok := repo.orders["order-002"]; !ok || updated.Status != domain.OrderApproved {
		t.Error("审核后订单状态应为 approved")
	}
}

// TestCancelOrderEmitsEvent 验证取消订单时正确发布 order.cancelled 事件
func TestCancelOrderEmitsEvent(t *testing.T) {
	repo := newMockOrderRepo()
	outboxStore := outbox.NewMemOutboxStore()
	svc := NewOrderAppService(repo).WithOutbox(outboxStore)

	order := &domain.SalesOrder{
		ID:             "order-003",
		TenantID:       "tenant-003",
		StoreID:        "store-001",
		PlatformOrderNo: "PO-2024-003",
		Status:          domain.OrderPending,
		IdempotencyKey:  "ik-003",
	}
	repo.orders[order.ID] = order

	if err := svc.CancelOrder(context.Background(), "order-003", "admin", "测试取消"); err != nil {
		t.Fatalf("取消订单失败: %v", err)
	}

	msgs := assertPendingMessages(t, outboxStore, 1)
	assertEventPayload(t, msgs[0], sharedEvents.EventOrderCancelled, "order-003")
	if msgs[0].TenantID != "tenant-003" {
		t.Errorf("期望租户ID tenant-003，实际 %s", msgs[0].TenantID)
	}
}

// TestEventsWithoutOutbox 验证未注入 outbox 时不发布事件（不报错）
func TestEventsWithoutOutbox(t *testing.T) {
	repo := newMockOrderRepo()
	svc := NewOrderAppService(repo)

	order := &domain.SalesOrder{
		ID:             "order-004",
		TenantID:       "tenant-004",
		StoreID:        "store-001",
		PlatformOrderNo: "PO-2024-004",
		Status:          domain.OrderPending,
		IdempotencyKey:  "ik-004",
	}
	if err := svc.CreateOrder(context.Background(), order); err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}
}
