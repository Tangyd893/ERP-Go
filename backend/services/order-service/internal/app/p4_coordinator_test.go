package app

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
)

// mockStockHandler 测试用库存锁定处理器
type mockStockHandler struct {
	locked           map[string][]string
	shouldFail       bool
	shouldFailOnCall int
	callCount        int
}

func newMockStockHandler() *mockStockHandler {
	return &mockStockHandler{locked: make(map[string][]string)}
}

func (h *mockStockHandler) LockStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) ([]string, error) {
	h.callCount++
	if h.shouldFail || (h.shouldFailOnCall > 0 && h.callCount == h.shouldFailOnCall) {
		return nil, fmt.Errorf("库存不足: 订单=%s", orderID)
	}
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
	outbounds  map[string]string
	shouldFail bool
}

func newMockOutboundCreator() *mockOutboundCreator {
	return &mockOutboundCreator{outbounds: make(map[string]string)}
}

func (h *mockOutboundCreator) CreateOutbound(ctx context.Context, tenantID, orderID, orderNo, warehouseID string, items []workflows.OrderItemData) (string, error) {
	if h.shouldFail {
		return "", fmt.Errorf("创建出库单失败: 仓库 %s 无可用人员", warehouseID)
	}
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

	data := workflows.OrderApprovedData{
		OrderID:     "order-001",
		TenantID:    "t-001",
		StoreID:     "st-001",
		OrderNo:     "AMZ-20260522-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
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

	if err := coordinator.HandleOrderApproved(context.Background(), "msg-001", payload); err != nil {
		t.Fatalf("重复处理同一消息应幂等返回: %v", err)
	}

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
			var lockData workflows.StockLockedData
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

	payload, err := outbox.NewEventPayload(events.EventOrderCancelled, workflows.OrderCancelledData{
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

func TestP4HandleOrderApprovedInsufficientStock(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	stockHandler := newMockStockHandler()
	stockHandler.shouldFail = true
	coordinator.SetStockHandler(stockHandler)
	coordinator.SetOutboundCreator(newMockOutboundCreator())

	data := workflows.OrderApprovedData{
		OrderID:     "order-100",
		TenantID:    "t-001",
		OrderNo:     "AMZ-INSUF-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-999", SKUCode: "OUTOFSTOCK", SKUName: "缺货SKU", Qty: 1000},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOrderApproved, data)

	err := coordinator.HandleOrderApproved(context.Background(), "msg-100", payload)
	if err == nil {
		t.Error("库存不足时应返回错误")
	}

	pending, _ := store.FetchPending(context.Background(), 10)
	if len(pending) > 0 {
		t.Errorf("库存不足时不应产生后续事件，实际有 %d 条", len(pending))
	}

	dup, _ := inbox.IsDuplicate(context.Background(), "msg-100")
	if dup {
		t.Error("库存不足时不应写入 inbox（未成功处理）")
	}
}

func TestP4HandleOrderApprovedDuplicateEvent(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(newMockStockHandler())
	coordinator.SetOutboundCreator(newMockOutboundCreator())

	data := workflows.OrderApprovedData{
		OrderID:     "order-200",
		TenantID:    "t-001",
		OrderNo:     "AMZ-DUP-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 1},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOrderApproved, data)

	if err := coordinator.HandleOrderApproved(context.Background(), "msg-200", payload); err != nil {
		t.Fatalf("首次处理失败: %v", err)
	}

	pending1, _ := store.FetchPending(context.Background(), 10)
	count1 := countOutboxByType(pending1, events.EventStockLocked)

	if err := coordinator.HandleOrderApproved(context.Background(), "msg-200", payload); err != nil {
		t.Fatalf("重复处理应幂等返回: %v", err)
	}

	pending2, _ := store.FetchPending(context.Background(), 10)
	count2 := countOutboxByType(pending2, events.EventStockLocked)

	if count2 != count1 {
		t.Errorf("重复事件不应产生新的库存锁定事件: 首次=%d, 重复后=%d", count1, count2)
	}
}

func TestP4HandleOrderApprovedOutboundFailed(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(newMockStockHandler())

	creator := newMockOutboundCreator()
	creator.shouldFail = true
	coordinator.SetOutboundCreator(creator)

	data := workflows.OrderApprovedData{
		OrderID:     "order-300",
		TenantID:    "t-001",
		OrderNo:     "AMZ-OBFAIL-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 1},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOrderApproved, data)

	err := coordinator.HandleOrderApproved(context.Background(), "msg-300", payload)
	if err == nil {
		t.Error("出库单创建失败时应返回错误")
	}

	pending, _ := store.FetchPending(context.Background(), 10)
	foundLock := false
	foundOutbound := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockLocked {
			foundLock = true
		}
		if msg.EventType == events.EventOutboundCreated {
			foundOutbound = true
		}
	}
	if !foundLock {
		t.Error("库存锁定应已执行，应产生 inventory.locked 事件")
	}
	if foundOutbound {
		t.Error("出库单创建失败时不应产生 outbound.created 事件")
	}

	dup, _ := inbox.IsDuplicate(context.Background(), "msg-300")
	if dup {
		t.Error("出库单创建失败时不应写入 inbox（待重试）")
	}
}

func countOutboxByType(msgs []*outbox.OutboxMessage, eventType string) int {
	n := 0
	for _, m := range msgs {
		if m.EventType == eventType {
			n++
		}
	}
	return n
}

type mockStockDeductHandler struct {
	deducted   map[string]bool
	shouldFail bool
}

func newMockStockDeductHandler() *mockStockDeductHandler {
	return &mockStockDeductHandler{deducted: make(map[string]bool)}
}

func (h *mockStockDeductHandler) DeductStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) error {
	if h.shouldFail {
		return fmt.Errorf("库存扣减失败: 订单=%s 部分SKU已过期", orderID)
	}
	h.deducted[orderID] = true
	return nil
}

type mockOrderStatusUpdater struct {
	statuses   map[string]string
	shouldFail bool
}

func newMockOrderStatusUpdater() *mockOrderStatusUpdater {
	return &mockOrderStatusUpdater{statuses: make(map[string]string)}
}

func (h *mockOrderStatusUpdater) UpdateOrderStatus(ctx context.Context, orderID string, status string, metadata map[string]interface{}) error {
	if h.shouldFail {
		return fmt.Errorf("订单状态更新失败: 订单=%s", orderID)
	}
	h.statuses[orderID] = status
	return nil
}

type mockCompensationStore struct {
	compensations []workflows.CompensationRecord
}

func newMockCompensationStore() *mockCompensationStore {
	return &mockCompensationStore{}
}

func (s *mockCompensationStore) CreateCompensation(ctx context.Context, orderID string, eventType string, payload []byte, reason string) error {
	s.compensations = append(s.compensations, workflows.CompensationRecord{
		ID:        fmt.Sprintf("comp-%d", len(s.compensations)+1),
		OrderID:   orderID,
		EventType: eventType,
		Payload:   payload,
		Reason:    reason,
		Status:    "pending",
		CreatedAt: time.Now(),
	})
	return nil
}

func (s *mockCompensationStore) FetchPendingCompensations(ctx context.Context, limit int) ([]workflows.CompensationRecord, error) {
	return s.compensations, nil
}

func (s *mockCompensationStore) MarkCompensationResolved(ctx context.Context, id string) error {
	for i, c := range s.compensations {
		if c.ID == id {
			s.compensations[i].Status = "resolved"
			break
		}
	}
	return nil
}

func TestP4HandleOutboundShipped(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockDeductHandler(newMockStockDeductHandler())
	coordinator.SetOrderStatusUpdater(newMockOrderStatusUpdater())
	coordinator.SetCompensationStore(newMockCompensationStore())

	data := workflows.OutboundShippedData{
		OutboundID:  "OB-001",
		OrderID:     "order-400",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "SF123456789",
		Carrier:     "顺丰速运",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 2},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOutboundShipped, data)

	if err := coordinator.HandleOutboundShipped(context.Background(), "msg-400", payload); err != nil {
		t.Fatalf("处理出库发货事件失败: %v", err)
	}

	if err := coordinator.HandleOutboundShipped(context.Background(), "msg-400", payload); err != nil {
		t.Fatalf("重复处理应幂等: %v", err)
	}

	pending, _ := store.FetchPending(context.Background(), 10)
	foundDeduct := false
	foundShipped := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockDeducted {
			foundDeduct = true
		}
		if msg.EventType == events.EventOrderShipped {
			foundShipped = true
		}
	}
	if !foundDeduct {
		t.Error("应产生 inventory.deducted 事件")
	}
	if !foundShipped {
		t.Error("应产生 order.shipped 事件")
	}
}

func TestP4HandleOutboundShippedDeductFailed(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	deductHandler := newMockStockDeductHandler()
	deductHandler.shouldFail = true
	coordinator.SetStockDeductHandler(deductHandler)
	coordinator.SetOrderStatusUpdater(newMockOrderStatusUpdater())

	compStore := newMockCompensationStore()
	coordinator.SetCompensationStore(compStore)

	data := workflows.OutboundShippedData{
		OutboundID:  "OB-002",
		OrderID:     "order-500",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "SF000",
		Carrier:     "顺丰速运",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-999", SKUCode: "EXPIRED", SKUName: "过期SKU", Qty: 10},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOutboundShipped, data)

	err := coordinator.HandleOutboundShipped(context.Background(), "msg-500", payload)
	if err == nil {
		t.Error("库存扣减失败时应返回错误")
	}

	if len(compStore.compensations) != 1 {
		t.Fatalf("库存扣减失败应写入 1 条补偿记录，实际: %d", len(compStore.compensations))
	}
	if compStore.compensations[0].OrderID != "order-500" {
		t.Errorf("补偿记录订单ID应为 order-500，实际: %s", compStore.compensations[0].OrderID)
	}
	if compStore.compensations[0].Status != "pending" {
		t.Errorf("补偿记录状态应为 pending，实际: %s", compStore.compensations[0].Status)
	}

	dup, _ := inbox.IsDuplicate(context.Background(), "msg-500")
	if dup {
		t.Error("库存扣减失败时不应写入 inbox（待重试）")
	}
}

func TestP4HandleOutboundShippedStatusUpdateFailed(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	coordinator.SetStockDeductHandler(newMockStockDeductHandler())
	statusUpdater := newMockOrderStatusUpdater()
	statusUpdater.shouldFail = true
	coordinator.SetOrderStatusUpdater(statusUpdater)

	compStore := newMockCompensationStore()
	coordinator.SetCompensationStore(compStore)

	data := workflows.OutboundShippedData{
		OutboundID:  "OB-003",
		OrderID:     "order-600",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "YD111",
		Carrier:     "韵达",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 1},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOutboundShipped, data)

	err := coordinator.HandleOutboundShipped(context.Background(), "msg-600", payload)
	if err == nil {
		t.Error("订单状态更新失败时应返回错误")
	}

	if len(compStore.compensations) != 1 {
		t.Fatalf("订单状态更新失败应写入 1 条补偿记录，实际: %d", len(compStore.compensations))
	}

	pending, _ := store.FetchPending(context.Background(), 10)
	foundDeduct := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockDeducted {
			foundDeduct = true
		}
	}
	if !foundDeduct {
		t.Error("库存扣减成功时应产生 inventory.deducted 事件")
	}
}
