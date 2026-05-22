package workflows

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

// mockStockHandler 测试用库存锁定处理器
type mockStockHandler struct {
	locked map[string][]string
}

func newMockStockHandler() *mockStockHandler {
	return &mockStockHandler{locked: make(map[string][]string)}
}

func (h *mockStockHandler) LockStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) ([]string, error) {
	keys := make([]string, 0)
	for sku, qty := range skuQtys {
		_ = qty
		key := "lock-" + orderID + "-" + sku
		keys = append(keys, key)
		h.locked[orderID] = append(h.locked[orderID], key)
	}
	return keys, nil
}

// mockOutboundCreator 测试用出库单创建处理器
type mockOutboundCreator struct {
	outbounds map[string]string
}

func newMockOutboundCreator() *mockOutboundCreator {
	return &mockOutboundCreator{outbounds: make(map[string]string)}
}

func (h *mockOutboundCreator) CreateOutbound(ctx context.Context, tenantID, orderID, orderNo, warehouseID string, items []OrderItemData) (string, error) {
	id := "OB-" + orderID
	h.outbounds[orderID] = id
	return id, nil
}

func TestP4HandleOrderApproved(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(newMockStockHandler())
	coordinator.SetOutboundCreator(newMockOutboundCreator())

	data := OrderApprovedData{
		OrderID:    "order-001",
		TenantID:   "t-001",
		StoreID:    "st-001",
		OrderNo:    "AMZ-20260522-001",
		WarehouseID: "wh-001",
		Items: []OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤经典款", Qty: 2},
			{SKUID: "sku-002", SKUCode: "MUG-001", SKUName: "马克杯", Qty: 1},
		},
	}

	payload, err := outbox.NewEventPayload(events.EventOrderApproved, data)
	if err != nil {
		t.Fatalf("构建事件载荷失败: %v", err)
	}

	if err := coordinator.HandleOrderApproved(context.Background(), "msg-001", payload); err != nil {
		t.Fatalf("处理订单审核事件失败: %v", err)
	}

	// 验证幂等性
	if err := coordinator.HandleOrderApproved(context.Background(), "msg-001", payload); err != nil {
		t.Fatalf("重复处理同一消息应幂等返回: %v", err)
	}

	// 验证 outbox 中产生了后续事件
	pending, err := store.FetchPending(context.Background(), 10)
	if err != nil {
		t.Fatalf("获取 pending 消息失败: %v", err)
	}

	foundLock := false
	foundOutbound := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockLocked {
			foundLock = true
			var ep outbox.EventPayload
			json.Unmarshal(msg.Payload, &ep)
			var lockData StockLockedData
			json.Unmarshal(ep.Data, &lockData)
			if lockData.OrderID != "order-001" {
				t.Errorf("锁定事件 orderID 应为 order-001，实际: %s", lockData.OrderID)
			}
		}
		if msg.EventType == events.EventOutboundCreated {
			foundOutbound = true
		}
	}

	if !foundLock {
		t.Error("应产生 inventory.locked 事件")
	}
	if !foundOutbound {
		t.Error("应产生 warehouse.outbound.created 事件")
	}
}

func TestP4HandleOrderCancelled(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	payload, err := outbox.NewEventPayload(events.EventOrderCancelled, OrderCancelledData{
		OrderID: "order-002",
		Reason:  "客户取消",
	})
	if err != nil {
		t.Fatalf("构建事件载荷失败: %v", err)
	}

	if err := coordinator.HandleOrderCancelled(context.Background(), "msg-002", payload); err != nil {
		t.Fatalf("处理订单取消事件失败: %v", err)
	}

	pending, _ := store.FetchPending(context.Background(), 10)
	found := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockReleased {
			found = true
		}
	}
	if !found {
		t.Error("应产生 inventory.released 事件")
	}
}

func TestNewEventPayload(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	data := testData{Name: "test", Value: 42}

	payload, err := outbox.NewEventPayload("test.event", data)
	if err != nil {
		t.Fatalf("构建事件载荷失败: %v", err)
	}

	var ep outbox.EventPayload
	if err := json.Unmarshal(payload, &ep); err != nil {
		t.Fatalf("解析事件载荷失败: %v", err)
	}
	if ep.EventType != "test.event" {
		t.Errorf("事件类型应为 test.event，实际: %s", ep.EventType)
	}

	var parsed testData
	if err := json.Unmarshal(ep.Data, &parsed); err != nil {
		t.Fatalf("解析数据载荷失败: %v", err)
	}
	if parsed.Value != 42 {
		t.Errorf("值应为 42，实际: %d", parsed.Value)
	}
}

func TestMemOutboxStore(t *testing.T) {
	store := outbox.NewMemOutboxStore()

	msg := &outbox.OutboxMessage{
		AggregateID:   "order-001",
		AggregateType: "SalesOrder",
		EventType:     "order.created",
		Payload:       []byte(`{"test":true}`),
		Status:        outbox.StatusPending,
		CreatedAt:     time.Now(),
	}
	if err := store.Save(context.Background(), msg); err != nil {
		t.Fatalf("保存消息失败: %v", err)
	}
	if msg.ID != 1 {
		t.Errorf("消息 ID 应为 1，实际: %d", msg.ID)
	}

	pending, err := store.FetchPending(context.Background(), 10)
	if err != nil {
		t.Fatalf("获取 pending 消息失败: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("应有 1 条 pending 消息，实际: %d", len(pending))
	}

	if err := store.MarkPublished(context.Background(), msg.ID); err != nil {
		t.Fatalf("标记已发布失败: %v", err)
	}

	pending, _ = store.FetchPending(context.Background(), 10)
	if len(pending) != 0 {
		t.Errorf("标记已发布后 pending 消息应为 0，实际: %d", len(pending))
	}
}

func TestMemInboxStore(t *testing.T) {
	store := outbox.NewMemInboxStore()

	dup, err := store.IsDuplicate(context.Background(), "msg-001")
	if err != nil {
		t.Fatalf("检查重复失败: %v", err)
	}
	if dup {
		t.Error("新消息不应为重复")
	}

	if err := store.Save(context.Background(), &outbox.InboxMessage{
		MessageID: "msg-001", EventType: "test.event",
		Payload: []byte(`{}`), ProcessedAt: time.Now(),
	}); err != nil {
		t.Fatalf("保存 inbox 消息失败: %v", err)
	}

	dup, _ = store.IsDuplicate(context.Background(), "msg-001")
	if !dup {
		t.Error("已保存的消息应标记为重复")
	}
}
