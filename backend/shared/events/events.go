package events

// 业务事件类型定义，用于 Outbox/Inbox 模式的事件驱动通信

const (
	// 订单事件
	EventOrderImported   = "order.imported"
	EventOrderApproved   = "order.approved"
	EventOrderCancelled  = "order.cancelled"
	EventOrderAbnormal   = "order.abnormal"
	EventOrderShipped    = "order.shipped"
	EventOrderDelivered  = "order.delivered"
	EventOrderCompleted  = "order.completed"

	// 库存事件
	EventStockLocked     = "inventory.locked"
	EventStockReleased   = "inventory.released"
	EventStockDeducted   = "inventory.deducted"
	EventStockIncreased  = "inventory.increased"

	// 出库事件
	EventOutboundCreated = "warehouse.outbound.created"
	EventOutboundPicked  = "warehouse.outbound.picked"
	EventOutboundChecked = "warehouse.outbound.checked"
	EventOutboundPacked  = "warehouse.outbound.packed"
	EventOutboundWeighed = "warehouse.outbound.weighed"
	EventOutboundShipped = "warehouse.outbound.shipped"

	// 物流事件
	EventShipmentCreated  = "transport.shipment.created"
	EventLabelGenerated   = "transport.label.generated"
	EventTrackingUpdated  = "transport.tracking.updated"

	// 财务事件
	EventSettlementImported = "finance.settlement.imported"
	EventProfitCalculated   = "finance.profit.calculated"
)
