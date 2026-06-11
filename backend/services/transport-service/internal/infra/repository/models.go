package repository

import "time"

type ShipmentModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	TenantID    string    `gorm:"column:tenant_id;index"`
	OrderID     string    `gorm:"column:order_id"`
	OutboundID  string    `gorm:"column:outbound_id"`
	CarrierCode string    `gorm:"column:carrier_code"`
	ServiceCode string    `gorm:"column:service_code"`
	TrackingNo  string    `gorm:"column:tracking_no"`
	LabelURL    string    `gorm:"column:label_url"`
	Status      string    `gorm:"column:status;index"`
	Weight      float64   `gorm:"column:weight"`
	ShippingCost float64  `gorm:"column:shipping_cost"`
	Currency    string    `gorm:"column:currency"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
func (ShipmentModel) TableName() string { return "shipments" }

type CarrierModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
func (CarrierModel) TableName() string { return "carriers" }

// CarrierServiceModel 物流产品持久化模型
type CarrierServiceModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	CarrierID   string `gorm:"column:carrier_id;index"`
	Name        string `gorm:"column:name"`
	Code        string `gorm:"column:code"`
	ServiceType string `gorm:"column:service_type"`
}
func (CarrierServiceModel) TableName() string { return "carrier_services" }

// ShippingRuleModel 物流规则持久化模型
type ShippingRuleModel struct {
	ID               string  `gorm:"column:id;primaryKey"`
	TenantID         string  `gorm:"column:tenant_id;index"`
	Name             string  `gorm:"column:name"`
	Priority         int     `gorm:"column:priority"`
	CountryCodes     string  `gorm:"column:country_codes"` // JSON array stored as string
	MinWeight        float64 `gorm:"column:min_weight"`
	MaxWeight        float64 `gorm:"column:max_weight"`
	CarrierServiceID string  `gorm:"column:carrier_service_id"`
}
func (ShippingRuleModel) TableName() string { return "shipping_rules" }
