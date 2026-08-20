package workflows

import (
	"context"
	"time"
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

// T-601: coordinator moved to order-service/internal/app (inbound to warehouse-service/internal/app)
