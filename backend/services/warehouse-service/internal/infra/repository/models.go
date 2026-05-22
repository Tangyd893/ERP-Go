package repository

import "time"

type OutboundOrderModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	TenantID    string    `gorm:"column:tenant_id;index"`
	OrderID     string    `gorm:"column:order_id"`
	OrderNo     string    `gorm:"column:order_no"`
	WarehouseID string    `gorm:"column:warehouse_id"`
	Status      string    `gorm:"column:status;index"`
	WaveID      string    `gorm:"column:wave_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
func (OutboundOrderModel) TableName() string { return "outbound_orders" }

type OutboundItemModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	OutboundID  string `gorm:"column:outbound_id;index"`
	SKUID       string `gorm:"column:sku_id"`
	SKUCode     string `gorm:"column:sku_code"`
	SKUName     string `gorm:"column:sku_name"`
	Quantity    int    `gorm:"column:quantity"`
	PickedQty   int    `gorm:"column:picked_quantity"`
	CheckedQty  int    `gorm:"column:checked_quantity"`
	LocationID  string `gorm:"column:location_id"`
}
func (OutboundItemModel) TableName() string { return "outbound_items" }

type WarehouseModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	Address   string    `gorm:"column:address"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
func (WarehouseModel) TableName() string { return "warehouses" }

type PickTaskModel struct {
	ID           string `gorm:"column:id;primaryKey"`
	WaveID       string `gorm:"column:wave_id;index"`
	OutboundID   string `gorm:"column:outbound_id"`
	SKUID        string `gorm:"column:sku_id"`
	SKUCode      string `gorm:"column:sku_code"`
	SKUName      string `gorm:"column:sku_name"`
	Quantity     int    `gorm:"column:quantity"`
	PickedQty    int    `gorm:"column:picked_quantity"`
	LocationCode string `gorm:"column:location_code"`
	Status       string `gorm:"column:status"`
	PickerID     string `gorm:"column:picker_id"`
}
func (PickTaskModel) TableName() string { return "pick_tasks" }
