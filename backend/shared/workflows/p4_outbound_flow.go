package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

const (
	errCheckIdempotency = "检查幂等失败: %w"
	errParsePayload     = "解析事件载荷失败: %w"
)

// OrderApprovedData 订单审核通过事件数据
type OrderApprovedData struct {
	OrderID    string `json:"order_id"`
	TenantID   string `json:"tenant_id"`
	StoreID    string `json:"store_id"`
	OrderNo    string `json:"order_no"`
	WarehouseID string `json:"warehouse_id"`
	Items      []OrderItemData `json:"items"`
}

// OrderItemData 订单明细数据
type OrderItemData struct {
	SKUID   string  `json:"sku_id"`
	SKUCode string  `json:"sku_code"`
	SKUName string  `json:"sku_name"`
	Qty     int     `json:"quantity"`
}

// OrderCancelledData 订单取消事件数据
type OrderCancelledData struct {
	OrderID  string `json:"order_id"`
	TenantID string `json:"tenant_id"`
	Reason   string `json:"reason"`
}

// OutboundShippedData 出库完成事件数据
type OutboundShippedData struct {
	OutboundID  string          `json:"outbound_id"`
	OrderID     string          `json:"order_id"`
	TenantID    string          `json:"tenant_id"`
	WarehouseID string          `json:"warehouse_id"`
	Items       []OrderItemData `json:"items"`
	TrackingNo  string          `json:"tracking_no"`
	Carrier     string          `json:"carrier"`
}

// StockLockedData 库存锁定完成事件数据
type StockLockedData struct {
	OrderID    string `json:"order_id"`
	WarehouseID string `json:"warehouse_id"`
	LockKeys   []string `json:"lock_keys"`
}

// OutboundCreatedData 出库单创建完成事件数据
type OutboundCreatedData struct {
	OutboundID string `json:"outbound_id"`
	OrderID    string `json:"order_id"`
	OrderNo    string `json:"order_no"`
	Items      []OrderItemData `json:"items"`
}

// StockLockHandler 库存锁定处理器接口
type StockLockHandler interface {
	LockStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) ([]string, error)
}

// OutboundCreator 出库单创建器接口
type OutboundCreator interface {
	CreateOutbound(ctx context.Context, tenantID, orderID, orderNo, warehouseID string, items []OrderItemData) (string, error)
}

// StockDeductHandler 库存扣减处理器接口
type StockDeductHandler interface {
	DeductStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) error
}

// OrderStatusUpdater 订单状态更新器接口
type OrderStatusUpdater interface {
	UpdateOrderStatus(ctx context.Context, orderID string, status string, metadata map[string]interface{}) error
}

// CompensationStore 补偿记录存储接口（人工补偿入口）
type CompensationStore interface {
	CreateCompensation(ctx context.Context, orderID string, eventType string, payload []byte, reason string) error
	FetchPendingCompensations(ctx context.Context, limit int) ([]CompensationRecord, error)
	MarkCompensationResolved(ctx context.Context, id string) error
}

// CompensationRecord 补偿记录
type CompensationRecord struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	EventType string    `json:"event_type"`
	Payload   []byte    `json:"payload"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// InboundReceivedData 采购入库事件数据
type InboundReceivedData struct {
	InboundID   string          `json:"inbound_id"`
	PurchaseID  string          `json:"purchase_id"`
	TenantID    string          `json:"tenant_id"`
	WarehouseID string          `json:"warehouse_id"`
	SupplierID  string          `json:"supplier_id"`
	Items       []OrderItemData `json:"items"`
}

// InboundHandler 采购入库处理器接口
type InboundHandler interface {
	ReceiveInbound(ctx context.Context, tenantID, purchaseID, warehouseID, supplierID string, items []OrderItemData) (string, error)
}

// P4OutboundFlowCoordinator P4 订单履约流程协调器
type P4OutboundFlowCoordinator struct {
	outbox             outbox.OutboxStore
	inbox              outbox.InboxStore
	stockHandler       StockLockHandler
	stockDeductHandler StockDeductHandler
	outboundCreator    OutboundCreator
	orderStatusUpdater OrderStatusUpdater
	compensationStore  CompensationStore
	inboundHandler     InboundHandler
}

func NewP4OutboundFlowCoordinator(outboxStore outbox.OutboxStore, inboxStore outbox.InboxStore) *P4OutboundFlowCoordinator {
	return &P4OutboundFlowCoordinator{
		outbox: outboxStore,
		inbox:  inboxStore,
	}
}

func (c *P4OutboundFlowCoordinator) SetStockHandler(h StockLockHandler) {
	c.stockHandler = h
}

func (c *P4OutboundFlowCoordinator) SetOutboundCreator(h OutboundCreator) {
	c.outboundCreator = h
}

func (c *P4OutboundFlowCoordinator) SetStockDeductHandler(h StockDeductHandler) {
	c.stockDeductHandler = h
}

func (c *P4OutboundFlowCoordinator) SetOrderStatusUpdater(h OrderStatusUpdater) {
	c.orderStatusUpdater = h
}

func (c *P4OutboundFlowCoordinator) SetCompensationStore(h CompensationStore) {
	c.compensationStore = h
}

func (c *P4OutboundFlowCoordinator) SetInboundHandler(h InboundHandler) {
	c.inboundHandler = h
}

// HandleOrderApproved 处理订单审核通过事件
// 流程: 订单审核 → 锁定库存 → 创建出库单
func (c *P4OutboundFlowCoordinator) HandleOrderApproved(ctx context.Context, messageID string, payload []byte) error {
	isDup, err := c.inbox.IsDuplicate(ctx, messageID)
	if err != nil {
		return fmt.Errorf(errCheckIdempotency, err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf(errParsePayload, err)
	}

	var data OrderApprovedData
	if err := json.Unmarshal(eventData.Data, &data); err != nil {
		return fmt.Errorf("解析订单数据失败: %w", err)
	}

	skuQtys := make(map[string]int)
	for _, item := range data.Items {
		skuQtys[item.SKUID] = item.Qty
	}

	lockKeys, err := c.stockHandler.LockStock(ctx, data.OrderID, data.WarehouseID, skuQtys)
	if err != nil {
		return fmt.Errorf("锁定库存失败: %w", err)
	}

	stockPayload, _ := outbox.NewEventPayload(events.EventStockLocked, StockLockedData{
		OrderID: data.OrderID, WarehouseID: data.WarehouseID, LockKeys: lockKeys,
	})
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: data.OrderID, AggregateType: "SalesOrder",
		TenantID: data.TenantID,
		EventType: events.EventStockLocked, Payload: stockPayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	outboundID, err := c.outboundCreator.CreateOutbound(ctx, data.TenantID, data.OrderID, data.OrderNo, data.WarehouseID, data.Items)
	if err != nil {
		return fmt.Errorf("创建出库单失败: %w", err)
	}

	outboundPayload, _ := outbox.NewEventPayload(events.EventOutboundCreated, OutboundCreatedData{
		OutboundID: outboundID, OrderID: data.OrderID, OrderNo: data.OrderNo, Items: data.Items,
	})
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: outboundID, AggregateType: "OutboundOrder",
		TenantID: data.TenantID,
		EventType: events.EventOutboundCreated, Payload: outboundPayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	_ = c.inbox.Save(ctx, &outbox.InboxMessage{
		MessageID: messageID, EventType: events.EventOrderApproved,
		Payload: payload, ProcessedAt: time.Now(),
	})

	return nil
}

// HandleOrderCancelled 处理订单取消事件
// 流程: 订单取消 → 释放库存
func (c *P4OutboundFlowCoordinator) HandleOrderCancelled(ctx context.Context, messageID string, payload []byte) error {
	isDup, err := c.inbox.IsDuplicate(ctx, messageID)
	if err != nil {
		return fmt.Errorf(errCheckIdempotency, err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf(errParsePayload, err)
	}

	var data OrderCancelledData
	if err := json.Unmarshal(eventData.Data, &data); err != nil {
		return fmt.Errorf("解析订单取消数据失败: %w", err)
	}

	releasePayload, _ := outbox.NewEventPayload(events.EventStockReleased, eventData.Data)
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: data.OrderID, AggregateType: "SalesOrder",
		TenantID: data.TenantID,
		EventType: events.EventStockReleased, Payload: releasePayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	_ = c.inbox.Save(ctx, &outbox.InboxMessage{
		MessageID: messageID, EventType: events.EventOrderCancelled,
		Payload: payload, ProcessedAt: time.Now(),
	})

	return nil
}

// HandleOutboundShipped 处理出库完成事件
// 流程: 出库完成 → 库存扣减 → 订单发货
func (c *P4OutboundFlowCoordinator) HandleOutboundShipped(ctx context.Context, messageID string, payload []byte) error {
	isDup, err := c.inbox.IsDuplicate(ctx, messageID)
	if err != nil {
		return fmt.Errorf(errCheckIdempotency, err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf(errParsePayload, err)
	}

	var data OutboundShippedData
	if err := json.Unmarshal(eventData.Data, &data); err != nil {
		return fmt.Errorf("解析出库数据失败: %w", err)
	}

	skuQtys := make(map[string]int)
	for _, item := range data.Items {
		skuQtys[item.SKUID] = item.Qty
	}

	if c.stockDeductHandler != nil {
		if err := c.stockDeductHandler.DeductStock(ctx, data.OrderID, data.WarehouseID, skuQtys); err != nil {
			c.createCompensation(ctx, data.OrderID, events.EventOutboundShipped, payload, fmt.Sprintf("库存扣减失败: %v", err))
			return fmt.Errorf("库存扣减失败: %w", err)
		}
	}

	deductPayload, _ := outbox.NewEventPayload(events.EventStockDeducted, map[string]interface{}{
		"order_id":     data.OrderID,
		"warehouse_id": data.WarehouseID,
		"items":        data.Items,
		"tracking_no":  data.TrackingNo,
	})
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: data.OrderID, AggregateType: "SalesOrder",
		TenantID: data.TenantID,
		EventType: events.EventStockDeducted, Payload: deductPayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	if c.orderStatusUpdater != nil {
		if err := c.orderStatusUpdater.UpdateOrderStatus(ctx, data.OrderID, "shipped", map[string]interface{}{
			"tracking_no": data.TrackingNo,
			"carrier":     data.Carrier,
		}); err != nil {
			c.createCompensation(ctx, data.OrderID, events.EventOutboundShipped, payload, fmt.Sprintf("订单状态更新失败: %v", err))
			return fmt.Errorf("订单状态更新失败: %w", err)
		}
	}

	shipPayload, _ := outbox.NewEventPayload(events.EventOrderShipped, map[string]interface{}{
		"order_id":    data.OrderID,
		"tracking_no": data.TrackingNo,
		"carrier":     data.Carrier,
	})
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: data.OrderID, AggregateType: "SalesOrder",
		TenantID: data.TenantID,
		EventType: events.EventOrderShipped, Payload: shipPayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	_ = c.inbox.Save(ctx, &outbox.InboxMessage{
		MessageID: messageID, EventType: events.EventOutboundShipped,
		Payload: payload, ProcessedAt: time.Now(),
	})

	return nil
}

// HandleInboundReceived 处理采购入库事件
// 流程: 入库接收 → 记录入库 → 增加库存 → 记录入库成本
func (c *P4OutboundFlowCoordinator) HandleInboundReceived(ctx context.Context, messageID string, payload []byte) error {
	isDup, err := c.inbox.IsDuplicate(ctx, messageID)
	if err != nil {
		return fmt.Errorf(errCheckIdempotency, err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf(errParsePayload, err)
	}

	var data InboundReceivedData
	if err := json.Unmarshal(eventData.Data, &data); err != nil {
		return fmt.Errorf("解析入库数据失败: %w", err)
	}

	if c.inboundHandler != nil {
		inboundID, err := c.inboundHandler.ReceiveInbound(ctx, data.TenantID, data.PurchaseID, data.WarehouseID, data.SupplierID, data.Items)
		if err != nil {
			return fmt.Errorf("创建入库记录失败: %w", err)
		}
		data.InboundID = inboundID
	}

	for _, item := range data.Items {
		incPayload, _ := outbox.NewEventPayload(events.EventStockIncreased, map[string]interface{}{
			"inbound_id":   data.InboundID,
			"purchase_id":  data.PurchaseID,
			"warehouse_id": data.WarehouseID,
			"sku_id":       item.SKUID,
			"sku_code":     item.SKUCode,
			"quantity":     item.Qty,
		})
		_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
			AggregateID:   data.InboundID,
			AggregateType: "InboundOrder",
			TenantID:      data.TenantID,
			EventType:     events.EventStockIncreased,
			Payload:       incPayload,
			Status:        outbox.StatusPending,
			CreatedAt:     time.Now(),
		})
	}

	settlePayload, _ := outbox.NewEventPayload(events.EventSettlementImported, data)
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID:   data.InboundID,
		AggregateType: "InboundOrder",
		TenantID:      data.TenantID,
		EventType:     events.EventSettlementImported,
		Payload:       settlePayload,
		Status:        outbox.StatusPending,
		CreatedAt:     time.Now(),
	})

	_ = c.inbox.Save(ctx, &outbox.InboxMessage{
		MessageID:   messageID,
		EventType:   events.EventSettlementImported,
		Payload:     payload,
		ProcessedAt: time.Now(),
	})

	return nil
}

func (c *P4OutboundFlowCoordinator) createCompensation(ctx context.Context, orderID, eventType string, payload []byte, reason string) {
	if c.compensationStore == nil {
		return
	}
	_ = c.compensationStore.CreateCompensation(ctx, orderID, eventType, payload, reason)
}
