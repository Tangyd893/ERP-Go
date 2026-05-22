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
