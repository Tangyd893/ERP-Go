package repository

import "time"

type SupplierModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	TenantID     string    `gorm:"column:tenant_id;index"`
	Name         string    `gorm:"column:name"`
	Code         string    `gorm:"column:code"`
	ContactName  string    `gorm:"column:contact_name"`
	ContactPhone string    `gorm:"column:contact_phone"`
	Email        string    `gorm:"column:email"`
	PaymentTerm  string    `gorm:"column:payment_term"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}
func (SupplierModel) TableName() string { return "suppliers" }

type PurchaseOrderModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	TenantID     string    `gorm:"column:tenant_id;index"`
	SupplierID   string    `gorm:"column:supplier_id"`
	SupplierName string    `gorm:"column:supplier_name"`
	OrderNo      string    `gorm:"column:order_no"`
	Status       string    `gorm:"column:status;index"`
	Currency     string    `gorm:"column:currency"`
	TotalAmount  float64   `gorm:"column:total_amount"`
	ExpectedDate time.Time `gorm:"column:expected_date"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}
func (PurchaseOrderModel) TableName() string { return "purchase_orders" }

type PurchaseItemModel struct {
	ID          string  `gorm:"column:id;primaryKey"`
	OrderID     string  `gorm:"column:order_id;index"`
	SKUID       string  `gorm:"column:sku_id"`
	SKUCode     string  `gorm:"column:sku_code"`
	SKUName     string  `gorm:"column:sku_name"`
	Quantity    int     `gorm:"column:quantity"`
	ReceivedQty int     `gorm:"column:received_quantity"`
	UnitPrice   float64 `gorm:"column:unit_price"`
	TotalPrice  float64 `gorm:"column:total_price"`
}
func (PurchaseItemModel) TableName() string { return "purchase_items" }

type InboundOrderModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	TenantID    string    `gorm:"column:tenant_id;index"`
	PurchaseID  string    `gorm:"column:purchase_id"`
	WarehouseID string    `gorm:"column:warehouse_id"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
func (InboundOrderModel) TableName() string { return "inbound_orders" }
