package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
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
	Reason   string `json:"reason"`
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

// P4OutboundFlowCoordinator P4 订单履约流程协调器
type P4OutboundFlowCoordinator struct {
	outbox          outbox.OutboxStore
	inbox           outbox.InboxStore
	stockHandler    StockLockHandler
	outboundCreator OutboundCreator
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

// HandleOrderApproved 处理订单审核通过事件
// 流程: 订单审核 → 锁定库存 → 创建出库单
func (c *P4OutboundFlowCoordinator) HandleOrderApproved(ctx context.Context, messageID string, payload []byte) error {
	isDup, err := c.inbox.IsDuplicate(ctx, messageID)
	if err != nil {
		return fmt.Errorf("检查幂等失败: %w", err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf("解析事件载荷失败: %w", err)
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
		return fmt.Errorf("检查幂等失败: %w", err)
	}
	if isDup {
		return nil
	}

	var eventData outbox.EventPayload
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf("解析事件载荷失败: %w", err)
	}

	releasePayload, _ := outbox.NewEventPayload(events.EventStockReleased, eventData.Data)
	_ = c.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID: "cancel", AggregateType: "SalesOrder",
		EventType: events.EventStockReleased, Payload: releasePayload,
		Status: outbox.StatusPending, CreatedAt: time.Now(),
	})

	_ = c.inbox.Save(ctx, &outbox.InboxMessage{
		MessageID: messageID, EventType: events.EventOrderCancelled,
		Payload: payload, ProcessedAt: time.Now(),
	})

	return nil
}
