package repository

import "time"

type SalesOrderModel struct {
	ID              string     `gorm:"column:id;primaryKey"`
	TenantID        string     `gorm:"column:tenant_id;index"`
	StoreID         string     `gorm:"column:store_id;index"`
	PlatformOrderNo string     `gorm:"column:platform_order_no"`
	OrderType       string     `gorm:"column:order_type"`
	OrderSource     string     `gorm:"column:order_source"`
	OrderStatus     string     `gorm:"column:order_status;index"`
	BuyerName       string     `gorm:"column:buyer_name"`
	BuyerEmail      string     `gorm:"column:buyer_email"`
	Currency        string     `gorm:"column:currency"`
	TotalAmount     float64    `gorm:"column:total_amount"`
	ShippingAmount  float64    `gorm:"column:shipping_amount"`
	DiscountAmount  float64    `gorm:"column:discount_amount"`
	ActualAmount    float64    `gorm:"column:actual_amount"`
	IdempotencyKey  string     `gorm:"column:idempotency_key;unique"`
	Remark          string     `gorm:"column:remark"`
	OrderedAt       *time.Time `gorm:"column:ordered_at"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (SalesOrderModel) TableName() string { return "sales_orders" }

type OrderItemModel struct {
	ID         string    `gorm:"column:id;primaryKey"`
	OrderID    string    `gorm:"column:order_id;index"`
	SKUID      string    `gorm:"column:sku_id"`
	SKUCode    string    `gorm:"column:sku_code"`
	SKUName    string    `gorm:"column:sku_name"`
	Quantity   int       `gorm:"column:quantity"`
	UnitPrice  float64   `gorm:"column:unit_price"`
	TotalPrice float64   `gorm:"column:total_price"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (OrderItemModel) TableName() string { return "order_items" }

type OrderAddressModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	OrderID      string    `gorm:"column:order_id;unique"`
	ContactName  string    `gorm:"column:contact_name"`
	Phone        string    `gorm:"column:phone"`
	Email        string    `gorm:"column:email"`
	Country      string    `gorm:"column:country"`
	State        string    `gorm:"column:state"`
	City         string    `gorm:"column:city"`
	District     string    `gorm:"column:district"`
	AddressLine1 string    `gorm:"column:address_line1"`
	AddressLine2 string    `gorm:"column:address_line2"`
	PostalCode   string    `gorm:"column:postal_code"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (OrderAddressModel) TableName() string { return "order_addresses" }

type OrderStatusLogModel struct {
	ID         string    `gorm:"column:id;primaryKey"`
	OrderID    string    `gorm:"column:order_id;index"`
	FromStatus string    `gorm:"column:from_status"`
	ToStatus   string    `gorm:"column:to_status"`
	Operator   string    `gorm:"column:operator"`
	Remark     string    `gorm:"column:remark"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (OrderStatusLogModel) TableName() string { return "order_status_logs" }
