package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/telemetry"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
	"go.opentelemetry.io/otel/attribute"
)

const (
	errCheckIdempotency = "检查幂等失败: %w"
	errParsePayload     = "解析事件载荷失败: %w"
)

// P4OutboundFlowCoordinator P4 订单履约流程协调器
type P4OutboundFlowCoordinator struct {
	outbox             outbox.OutboxStore
	inbox              outbox.InboxStore
	stockHandler       workflows.StockLockHandler
	stockDeductHandler workflows.StockDeductHandler
	outboundCreator    workflows.OutboundCreator
	orderStatusUpdater workflows.OrderStatusUpdater
	compensationStore  workflows.CompensationStore
}

func NewP4OutboundFlowCoordinator(outboxStore outbox.OutboxStore, inboxStore outbox.InboxStore) *P4OutboundFlowCoordinator {
	return &P4OutboundFlowCoordinator{
		outbox: outboxStore,
		inbox:  inboxStore,
	}
}

func (c *P4OutboundFlowCoordinator) SetStockHandler(h workflows.StockLockHandler) {
	c.stockHandler = h
}

func (c *P4OutboundFlowCoordinator) SetOutboundCreator(h workflows.OutboundCreator) {
	c.outboundCreator = h
}

func (c *P4OutboundFlowCoordinator) SetStockDeductHandler(h workflows.StockDeductHandler) {
	c.stockDeductHandler = h
}

func (c *P4OutboundFlowCoordinator) SetOrderStatusUpdater(h workflows.OrderStatusUpdater) {
	c.orderStatusUpdater = h
}

func (c *P4OutboundFlowCoordinator) SetCompensationStore(h workflows.CompensationStore) {
	c.compensationStore = h
}

// HandleOrderApproved 处理订单审核通过事件
// 流程: 订单审核 → 锁定库存 → 创建出库单
func (c *P4OutboundFlowCoordinator) HandleOrderApproved(ctx context.Context, messageID string, payload []byte) error {
	ctx, end := telemetry.StartSpan(ctx, "P4.HandleOrderApproved",
		attribute.String("messaging.message_id", messageID),
	)
	defer end()

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

	var data workflows.OrderApprovedData
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

	stockPayload, _ := outbox.NewEventPayload(events.EventStockLocked, workflows.StockLockedData{
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

	outboundPayload, _ := outbox.NewEventPayload(events.EventOutboundCreated, workflows.OutboundCreatedData{
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
	ctx, end := telemetry.StartSpan(ctx, "P4.HandleOrderCancelled",
		attribute.String("messaging.message_id", messageID),
	)
	defer end()

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

	var data workflows.OrderCancelledData
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
	ctx, end := telemetry.StartSpan(ctx, "P4.HandleOutboundShipped",
		attribute.String("messaging.message_id", messageID),
	)
	defer end()

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

	var data workflows.OutboundShippedData
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
			"tenant_id":   data.TenantID,
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

func (c *P4OutboundFlowCoordinator) createCompensation(ctx context.Context, orderID, eventType string, payload []byte, reason string) {
	if c.compensationStore == nil {
		return
	}
	_ = c.compensationStore.CreateCompensation(ctx, orderID, eventType, payload, reason)
}
