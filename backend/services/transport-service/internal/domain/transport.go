package domain

import "time"

// ShipmentStatus 发运状态
type ShipmentStatus string

const (
	ShipmentPending    ShipmentStatus = "pending"
	ShipmentLabeled    ShipmentStatus = "labeled"
	ShipmentShipped    ShipmentStatus = "shipped"
	ShipmentInTransit  ShipmentStatus = "in_transit"
	ShipmentDelivered  ShipmentStatus = "delivered"
	ShipmentCancelled  ShipmentStatus = "cancelled"
	ShipmentFailed     ShipmentStatus = "failed"
)

// Carrier 物流商
type Carrier struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CarrierService 物流产品
type CarrierService struct {
	ID         string `json:"id"`
	CarrierID  string `json:"carrier_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	ServiceType string `json:"service_type"` // express, standard, economy
}

// ShippingRule 物流规则
type ShippingRule struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	Name         string `json:"name"`
	Priority     int    `json:"priority"`
	CountryCodes []string `json:"country_codes"`
	MinWeight    float64 `json:"min_weight"`
	MaxWeight    float64 `json:"max_weight"`
	CarrierServiceID string `json:"carrier_service_id"`
}

// Shipment 发运单聚合根
type Shipment struct {
	ID              string          `json:"id"`
	TenantID        string          `json:"tenant_id"`
	OrderID         string          `json:"order_id"`
	OutboundID      string          `json:"outbound_id"`
	CarrierCode     string          `json:"carrier_code"`
	ServiceCode     string          `json:"service_code"`
	TrackingNo      string          `json:"tracking_no"`
	LabelURL        string          `json:"label_url"`
	Status          ShipmentStatus  `json:"status"`
	Weight          float64         `json:"weight"`
	ShippingCost    float64         `json:"shipping_cost"`
	Currency        string          `json:"currency"`
	Packages        []*PackageInfo  `json:"packages"`
	TrackingRecords []*TrackingRecord `json:"tracking_records,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// PackageInfo 包裹信息
type PackageInfo struct {
	ID         string  `json:"id"`
	TrackingNo string  `json:"tracking_no"`
	Weight     float64 `json:"weight"`
	Length     float64 `json:"length"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
}

// TrackingRecord 物流轨迹
type TrackingRecord struct {
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	RecordedAt  time.Time `json:"recorded_at"`
}
