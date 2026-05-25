package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
)

// TestP4HTTPFulfillmentChain 验证 HTTP 适配器驱动的订单审核→出库→发货链路（内存协调器）
func TestP4HTTPFulfillmentChain(t *testing.T) {
	var lockedOrderID string
	var createdOutboundID string
	var deductedOrderID string

	inventorySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/inventory/lock":
			lockedOrderID = "ok"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 0})
		case "/api/v1/inventory/deduct-by-order":
			var body struct {
				OrderID string `json:"order_id"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			deductedOrderID = body.OrderID
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": map[string]bool{"deducted": true}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer inventorySrv.Close()

	warehouseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/warehouse/outbounds" && r.Method == http.MethodPost {
			createdOutboundID = "OB-TEST-001"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 0,
				"data": map[string]string{"id": createdOutboundID},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer warehouseSrv.Close()

	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := workflows.NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(workflows.NewHTTPStockLockAdapter(inventorySrv.URL))
	coordinator.SetStockDeductHandler(workflows.NewHTTPStockDeductAdapter(inventorySrv.URL))
	coordinator.SetOutboundCreator(workflows.NewHTTPOutboundCreatorAdapter(warehouseSrv.URL))

	orderID := "order-e2e-001"
	approvedPayload, _ := outbox.NewEventPayload(events.EventOrderApproved, workflows.OrderApprovedData{
		OrderID: orderID, TenantID: "default", StoreID: "st-1", OrderNo: "SO-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Qty: 2},
		},
	})
	if err := coordinator.HandleOrderApproved(context.Background(), "msg-approve-1", approvedPayload); err != nil {
		t.Fatalf("HandleOrderApproved: %v", err)
	}
	if lockedOrderID == "" {
		t.Fatal("应调用库存锁定")
	}
	if createdOutboundID == "" {
		t.Fatal("应创建出库单")
	}

	statusUpdated := false
	coordinator.SetOrderStatusUpdater(&mockStatusUpdater{
		onUpdate: func(orderID, status string) {
			if orderID == "order-e2e-001" && status == "shipped" {
				statusUpdated = true
			}
		},
	})

	shipPayload, _ := outbox.NewEventPayload(events.EventOutboundShipped, workflows.OutboundShippedData{
		OutboundID: createdOutboundID, OrderID: orderID, TenantID: "default",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{{SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Qty: 2}},
		TrackingNo: "TN-001", Carrier: "YTO",
	})
	if err := coordinator.HandleOutboundShipped(context.Background(), "ship-OB-TEST-001", shipPayload); err != nil {
		t.Fatalf("HandleOutboundShipped: %v", err)
	}
	if deductedOrderID != orderID {
		t.Fatalf("应扣减订单 %s 库存，实际 %s", orderID, deductedOrderID)
	}
	if !statusUpdated {
		t.Fatal("应更新订单为 shipped")
	}
}

// TestP4LockStockFailure 验证锁库失败时不会创建出库单，Inbox 未写入（事务回滚语义）
func TestP4LockStockFailure(t *testing.T) {
	var createdOutboundID string
	var deductedOrderID string

	inventorySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/inventory/lock":
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 1, "message": "库存不足"})
		case "/api/v1/inventory/deduct-by-order":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": map[string]bool{"deducted": true}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer inventorySrv.Close()

	warehouseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/warehouse/outbounds" && r.Method == http.MethodPost {
			createdOutboundID = "OB-SHOULD-NOT-CREATE"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 0,
				"data": map[string]string{"id": createdOutboundID},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer warehouseSrv.Close()

	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := workflows.NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(workflows.NewHTTPStockLockAdapter(inventorySrv.URL))
	coordinator.SetStockDeductHandler(workflows.NewHTTPStockDeductAdapter(inventorySrv.URL))
	coordinator.SetOutboundCreator(workflows.NewHTTPOutboundCreatorAdapter(warehouseSrv.URL))

	orderID := "order-lock-fail-001"
	approvedPayload, _ := outbox.NewEventPayload(events.EventOrderApproved, workflows.OrderApprovedData{
		OrderID: orderID, TenantID: "default", StoreID: "st-1", OrderNo: "SO-FAIL-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Qty: 2},
		},
	})

	// 锁库失败应返回错误且不创建出库单
	err := coordinator.HandleOrderApproved(context.Background(), "msg-approve-fail", approvedPayload)
	if err == nil {
		t.Fatal("锁库失败应返回错误")
	}
	if createdOutboundID != "" {
		t.Fatal("锁库失败后不应创建出库单")
	}

	// inbox 未写入，允许重试
	isDup, _ := inbox.IsDuplicate(context.Background(), "msg-approve-fail")
	if isDup {
		t.Fatal("锁库失败时不应写入 inbox")
	}

	// 确保没有因错误导致后续操作
	_ = deductedOrderID
}

// TestP4DuplicateOrderApproved 验证同一 messageID 两次调用仅处理一次（幂等）
func TestP4DuplicateOrderApproved(t *testing.T) {
	var lockCallCount int
	var outboundCreatedCount int

	inventorySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/inventory/lock" {
			lockCallCount++
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 0})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer inventorySrv.Close()

	warehouseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/warehouse/outbounds" && r.Method == http.MethodPost {
			outboundCreatedCount++
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 0,
				"data": map[string]string{"id": "OB-DUP-TEST"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer warehouseSrv.Close()

	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := workflows.NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockHandler(workflows.NewHTTPStockLockAdapter(inventorySrv.URL))
	coordinator.SetOutboundCreator(workflows.NewHTTPOutboundCreatorAdapter(warehouseSrv.URL))

	orderID := "order-dup-001"
	approvedPayload, _ := outbox.NewEventPayload(events.EventOrderApproved, workflows.OrderApprovedData{
		OrderID: orderID, TenantID: "default", StoreID: "st-1", OrderNo: "SO-DUP-001",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Qty: 2},
		},
	})

	// 第一次处理成功
	if err := coordinator.HandleOrderApproved(context.Background(), "msg-dup-1", approvedPayload); err != nil {
		t.Fatalf("第一次处理失败: %v", err)
	}
	if lockCallCount != 1 {
		t.Fatalf("第一次应调用一次锁库，实际 %d", lockCallCount)
	}
	if outboundCreatedCount != 1 {
		t.Fatalf("第一次应创建一次出库单，实际 %d", outboundCreatedCount)
	}

	// 第二次使用相同 messageID 应幂等跳过
	if err := coordinator.HandleOrderApproved(context.Background(), "msg-dup-1", approvedPayload); err != nil {
		t.Fatalf("第二次处理（幂等）不应报错: %v", err)
	}
	if lockCallCount != 1 {
		t.Fatalf("幂等调用不应再调锁库，实际 %d", lockCallCount)
	}
	if outboundCreatedCount != 1 {
		t.Fatalf("幂等调用不应再创出库单，实际 %d", outboundCreatedCount)
	}
}

// TestP4OutboundShippedViaConsumer 模拟 RabbitMQ Consumer handler 调 HandleOutboundShipped
func TestP4OutboundShippedViaConsumer(t *testing.T) {
	var deductedOrderID string
	var statusUpdatedTo string

	inventorySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/inventory/deduct-by-order" {
			var body struct {
				OrderID string `json:"order_id"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			deductedOrderID = body.OrderID
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": map[string]bool{"deducted": true}})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer inventorySrv.Close()

	store := outbox.NewMemOutboxStore()
	inbox := outbox.NewMemInboxStore()
	coordinator := workflows.NewP4OutboundFlowCoordinator(store, inbox)
	coordinator.SetStockDeductHandler(workflows.NewHTTPStockDeductAdapter(inventorySrv.URL))
	coordinator.SetOrderStatusUpdater(&mockStatusUpdater{
		onUpdate: func(orderID, status string) {
			statusUpdatedTo = status
		},
	})

	// 模拟 Consumer 接收到的消息载荷
	outboundID := "OB-CONSUMER-001"
	orderID := "order-consumer-001"
	shipPayload, _ := outbox.NewEventPayload(events.EventOutboundShipped, workflows.OutboundShippedData{
		OutboundID: outboundID, OrderID: orderID, TenantID: "default",
		WarehouseID: "wh-001",
		Items: []workflows.OrderItemData{
			{SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Qty: 3},
		},
		TrackingNo: "SF123", Carrier: "SF",
	})

	// 模拟 Consumer handler 调用
	consumerMessageID := "ship-" + outboundID
	if err := coordinator.HandleOutboundShipped(context.Background(), consumerMessageID, shipPayload); err != nil {
		t.Fatalf("Consumer HandleOutboundShipped 失败: %v", err)
	}

	if deductedOrderID != orderID {
		t.Fatalf("应扣减订单 %s 库存，实际 %s", orderID, deductedOrderID)
	}
	if statusUpdatedTo != "shipped" {
		t.Fatalf("应更新订单为 shipped，实际 %s", statusUpdatedTo)
	}

	// 幂等验证：重复消费同一消息不重复处理
	deductedOrderID = ""
	statusUpdatedTo = ""
	if err := coordinator.HandleOutboundShipped(context.Background(), consumerMessageID, shipPayload); err != nil {
		t.Fatalf("幂等消费不应报错: %v", err)
	}
	if deductedOrderID != "" {
		t.Fatal("幂等消费不应再次扣减")
	}
	if statusUpdatedTo != "" {
		t.Fatal("幂等消费不应再次更新状态")
	}
}
type mockStatusUpdater struct {
	onUpdate func(orderID, status string)
}

func (m *mockStatusUpdater) UpdateOrderStatus(_ context.Context, orderID, status string, _ map[string]interface{}) error {
	if m.onUpdate != nil {
		m.onUpdate(orderID, status)
	}
	return nil
}
