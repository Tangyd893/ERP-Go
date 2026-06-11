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

// WaveModel 波次持久化模型
type WaveModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	WarehouseID string    `gorm:"column:warehouse_id;index"`
	Name        string    `gorm:"column:name"`
	Status      string    `gorm:"column:status;index"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
func (WaveModel) TableName() string { return "waves" }

// WaveOutboundModel 波次-出库单关联
type WaveOutboundModel struct {
	ID         string `gorm:"column:id;primaryKey"`
	WaveID     string `gorm:"column:wave_id;index"`
	OutboundID string `gorm:"column:outbound_id;uniqueIndex"`
}
func (WaveOutboundModel) TableName() string { return "wave_outbounds" }

// PackageModel 包裹持久化模型
type PackageModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	OutboundID  string    `gorm:"column:outbound_id;index"`
	TrackingNo  string    `gorm:"column:tracking_no"`
	CarrierCode string    `gorm:"column:carrier_code"`
	Weight      float64   `gorm:"column:weight"`
	Length      float64   `gorm:"column:length"`
	Width       float64   `gorm:"column:width"`
	Height      float64   `gorm:"column:height"`
	LabelURL    string    `gorm:"column:label_url"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
func (PackageModel) TableName() string { return "packages" }
