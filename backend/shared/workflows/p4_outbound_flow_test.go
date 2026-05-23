package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

// mockStockHandler 测试用库存锁定处理器
type mockStockHandler struct {
	locked      map[string][]string
	shouldFail  bool
	shouldFailOnCall int
	callCount   int
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

func (h *mockOutboundCreator) CreateOutbound(ctx context.Context, tenantID, orderID, orderNo, warehouseID string, items []OrderItemData) (string, error) {
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

// TestP4HandleOrderApprovedInsufficientStock 验证库存不足时锁库失败，不创建出库单
func TestP4HandleOrderApprovedInsufficientStock(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	stockHandler := newMockStockHandler()
	stockHandler.shouldFail = true
	coordinator.SetStockHandler(stockHandler)
	coordinator.SetOutboundCreator(newMockOutboundCreator())

	data := OrderApprovedData{
		OrderID:     "order-100",
		TenantID:    "t-001",
		OrderNo:     "AMZ-INSUF-001",
		WarehouseID: "wh-001",
		Items: []OrderItemData{
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

// TestP4HandleOrderApprovedDuplicateEvent 验证同一事件重复投递时幂等处理
func TestP4HandleOrderApprovedDuplicateEvent(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(newMockStockHandler())
	coordinator.SetOutboundCreator(newMockOutboundCreator())

	data := OrderApprovedData{
		OrderID:     "order-200",
		TenantID:    "t-001",
		OrderNo:     "AMZ-DUP-001",
		WarehouseID: "wh-001",
		Items: []OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 1},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOrderApproved, data)

	if err := coordinator.HandleOrderApproved(context.Background(), "msg-200", payload); err != nil {
		t.Fatalf("首次处理失败: %v", err)
	}

	pending1, _ := store.FetchPending(context.Background(), 10)
	count1 := countOutboxByType(pending1, events.EventStockLocked)

	// 重复投递同一消息（相同 messageID）
	if err := coordinator.HandleOrderApproved(context.Background(), "msg-200", payload); err != nil {
		t.Fatalf("重复处理应幂等返回: %v", err)
	}

	pending2, _ := store.FetchPending(context.Background(), 10)
	count2 := countOutboxByType(pending2, events.EventStockLocked)

	if count2 != count1 {
		t.Errorf("重复事件不应产生新的库存锁定事件: 首次=%d, 重复后=%d", count1, count2)
	}
}

// TestP4HandleOrderApprovedOutboundFailed 验证出库单创建失败时不写 inbox
func TestP4HandleOrderApprovedOutboundFailed(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(newMockStockHandler())

	creator := newMockOutboundCreator()
	creator.shouldFail = true
	coordinator.SetOutboundCreator(creator)

	data := OrderApprovedData{
		OrderID:     "order-300",
		TenantID:    "t-001",
		OrderNo:     "AMZ-OBFAIL-001",
		WarehouseID: "wh-001",
		Items: []OrderItemData{
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

// mockStockDeductHandler 测试用库存扣减处理器
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

// mockOrderStatusUpdater 测试用订单状态更新器
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

// mockCompensationStore 测试用补偿存储
type mockCompensationStore struct {
	compensations []CompensationRecord
}

func newMockCompensationStore() *mockCompensationStore {
	return &mockCompensationStore{}
}

func (s *mockCompensationStore) CreateCompensation(ctx context.Context, orderID string, eventType string, payload []byte, reason string) error {
	s.compensations = append(s.compensations, CompensationRecord{
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

func (s *mockCompensationStore) FetchPendingCompensations(ctx context.Context, limit int) ([]CompensationRecord, error) {
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

// TestP4HandleOutboundShipped 验证出库完成→库存扣减→订单发货完整流程
func TestP4HandleOutboundShipped(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockDeductHandler(newMockStockDeductHandler())
	coordinator.SetOrderStatusUpdater(newMockOrderStatusUpdater())
	coordinator.SetCompensationStore(newMockCompensationStore())

	data := OutboundShippedData{
		OutboundID:  "OB-001",
		OrderID:     "order-400",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "SF123456789",
		Carrier:     "顺丰速运",
		Items: []OrderItemData{
			{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤", Qty: 2},
		},
	}
	payload, _ := outbox.NewEventPayload(events.EventOutboundShipped, data)

	if err := coordinator.HandleOutboundShipped(context.Background(), "msg-400", payload); err != nil {
		t.Fatalf("处理出库发货事件失败: %v", err)
	}

	// 验证幂等
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

// TestP4HandleOutboundShippedDeductFailed 验证库存扣减失败时写入补偿记录
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

	data := OutboundShippedData{
		OutboundID:  "OB-002",
		OrderID:     "order-500",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "SF000",
		Carrier:     "顺丰速运",
		Items: []OrderItemData{
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

// TestP4HandleOutboundShippedStatusUpdateFailed 验证订单状态更新失败时写入补偿
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

	data := OutboundShippedData{
		OutboundID:  "OB-003",
		OrderID:     "order-600",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		TrackingNo:  "YD111",
		Carrier:     "韵达",
		Items: []OrderItemData{
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

	// 库存扣减事件仍应产生（因为扣减成功了）
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

// mockInboundHandler 测试用采购入库处理器
type mockInboundHandler struct {
	inbounds   map[string]string
	shouldFail bool
}

func newMockInboundHandler() *mockInboundHandler {
	return &mockInboundHandler{inbounds: make(map[string]string)}
}

func (h *mockInboundHandler) ReceiveInbound(ctx context.Context, tenantID, purchaseID, warehouseID, supplierID string, items []OrderItemData) (string, error) {
	if h.shouldFail {
		return "", fmt.Errorf("创建入库记录失败: 仓库 %s 不可用", warehouseID)
	}
	id := "IB-" + purchaseID
	h.inbounds[purchaseID] = id
	return id, nil
}

// TestP4HandleInboundReceived 验证采购入库完整流程及幂等
func TestP4HandleInboundReceived(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetInboundHandler(newMockInboundHandler())

	data := InboundReceivedData{
		InboundID:   "IB-001",
		PurchaseID:  "PO-001",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		SupplierID:  "SUP-001",
		Items: []OrderItemData{
			{SKUID: "sku-001", SKUCode: "MAT-001", SKUName: "原材料A", Qty: 100},
			{SKUID: "sku-002", SKUCode: "MAT-002", SKUName: "原材料B", Qty: 50},
		},
	}

	payload, err := outbox.NewEventPayload(events.EventSettlementImported, data)
	if err != nil {
		t.Fatalf("构建事件载荷失败: %v", err)
	}

	if err := coordinator.HandleInboundReceived(context.Background(), "msg-700", payload); err != nil {
		t.Fatalf("处理采购入库事件失败: %v", err)
	}

	// 验证幂等性
	if err := coordinator.HandleInboundReceived(context.Background(), "msg-700", payload); err != nil {
		t.Fatalf("重复处理同一消息应幂等返回: %v", err)
	}

	pending, err := store.FetchPending(context.Background(), 20)
	if err != nil {
		t.Fatalf("获取 pending 消息失败: %v", err)
	}

	foundIncreased := 0
	foundSettlement := false
	for _, msg := range pending {
		if msg.EventType == events.EventStockIncreased {
			foundIncreased++
		}
		if msg.EventType == events.EventSettlementImported {
			foundSettlement = true
		}
	}

	if foundIncreased != 2 {
		t.Errorf("应为 2 个 SKU 各产生 inventory.increased 事件，实际: %d", foundIncreased)
	}
	if !foundSettlement {
		t.Error("应产生 finance.settlement.imported 事件")
	}
}

// TestP4HandleInboundReceivedInboundFailed 验证入库处理失败时不写 inbox
func TestP4HandleInboundReceivedInboundFailed(t *testing.T) {
	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := NewP4OutboundFlowCoordinator(store, inbox)

	handler := newMockInboundHandler()
	handler.shouldFail = true
	coordinator.SetInboundHandler(handler)

	data := InboundReceivedData{
		PurchaseID:  "PO-002",
		TenantID:    "t-001",
		WarehouseID: "wh-001",
		SupplierID:  "SUP-001",
		Items: []OrderItemData{
			{SKUID: "sku-001", SKUCode: "MAT-001", SKUName: "原材料A", Qty: 100},
		},
	}

	payload, _ := outbox.NewEventPayload(events.EventSettlementImported, data)

	err := coordinator.HandleInboundReceived(context.Background(), "msg-800", payload)
	if err == nil {
		t.Error("入库处理失败时应返回错误")
	}

	dup, _ := inbox.IsDuplicate(context.Background(), "msg-800")
	if dup {
		t.Error("入库处理失败时不应写入 inbox（待重试）")
	}
}
