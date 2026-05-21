package domain

import "time"

// OutboundStatus 出库单状态
type OutboundStatus string

const (
	OutboundCreated  OutboundStatus = "created"
	OutboundWaved    OutboundStatus = "waved"
	OutboundPicking  OutboundStatus = "picking"
	OutboundPicked   OutboundStatus = "picked"
	OutboundChecking OutboundStatus = "checking"
	OutboundChecked  OutboundStatus = "checked"
	OutboundPacking  OutboundStatus = "packing"
	OutboundPacked   OutboundStatus = "packed"
	OutboundWeighed  OutboundStatus = "weighed"
	OutboundShipped  OutboundStatus = "shipped"
	OutboundAbnormal OutboundStatus = "abnormal"
)

// OutboundOrder 出库单聚合根
type OutboundOrder struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id"`
	OrderID     string          `json:"order_id"`
	OrderNo     string          `json:"order_no"`
	WarehouseID string          `json:"warehouse_id"`
	Status      OutboundStatus  `json:"status"`
	WaveID      string          `json:"wave_id"`
	Items       []*OutboundItem `json:"items"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// OutboundItem 出库明细
type OutboundItem struct {
	ID        string `json:"id"`
	SKUID     string `json:"sku_id"`
	SKUCode   string `json:"sku_code"`
	SKUName   string `json:"sku_name"`
	Quantity  int    `json:"quantity"`
	PickedQty int    `json:"picked_quantity"`
	CheckedQty int   `json:"checked_quantity"`
	LocationID string  `json:"location_id"`
}

// Warehouse 仓库实体
type Warehouse struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Location 库位
type Location struct {
	ID          string `json:"id"`
	WarehouseID string `json:"warehouse_id"`
	ZoneID      string `json:"zone_id"`
	Code        string `json:"code"`
	Barcode     string `json:"barcode"`
	Status      string `json:"status"`
}

// Zone 库区
type Zone struct {
	ID          string `json:"id"`
	WarehouseID string `json:"warehouse_id"`
	Name        string `json:"name"`
	ZoneType    string `json:"zone_type"` // pick, reserve, receive, return
}

// Wave 波次（拣货批次）
type Wave struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouse_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"` // created, picking, completed
	OutboundIDs []string  `json:"outbound_ids"`
	CreatedAt   time.Time `json:"created_at"`
}

// PickTask 拣货任务
type PickTask struct {
	ID         string `json:"id"`
	WaveID     string `json:"wave_id"`
	OutboundID string `json:"outbound_id"`
	SKUID      string `json:"sku_id"`
	SKUCode    string `json:"sku_code"`
	SKUName    string `json:"sku_name"`
	Quantity   int    `json:"quantity"`
	PickedQty  int    `json:"picked_quantity"`
	LocationCode string `json:"location_code"`
	Status     string `json:"status"` // pending, picking, picked, check_pending, checked
	PickerID   string `json:"picker_id"`
}

// Package 包裹
type Package struct {
	ID          string    `json:"id"`
	OutboundID  string    `json:"outbound_id"`
	TrackingNo  string    `json:"tracking_no"`
	CarrierCode string    `json:"carrier_code"`
	Weight      float64   `json:"weight"`
	Length      float64   `json:"length"`
	Width       float64   `json:"width"`
	Height      float64   `json:"height"`
	LabelURL    string    `json:"label_url"`
	CreatedAt   time.Time `json:"created_at"`
}
