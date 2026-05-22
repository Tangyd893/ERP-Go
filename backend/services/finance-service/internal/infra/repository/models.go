package repository

import "time"

type SettlementBillModel struct {
	ID               string    `gorm:"column:id;primaryKey"`
	TenantID         string    `gorm:"column:tenant_id;index"`
	StoreID          string    `gorm:"column:store_id"`
	PlatformCode     string    `gorm:"column:platform_code"`
	SettlementPeriod string    `gorm:"column:settlement_period"`
	Currency         string    `gorm:"column:currency"`
	TotalSales       float64   `gorm:"column:total_sales"`
	TotalRefunds     float64   `gorm:"column:total_refunds"`
	Commission       float64   `gorm:"column:commission"`
	FbaFee           float64   `gorm:"column:fba_fee"`
	OtherFee         float64   `gorm:"column:other_fee"`
	NetAmount        float64   `gorm:"column:net_amount"`
	Status           string    `gorm:"column:status"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}
func (SettlementBillModel) TableName() string { return "settlement_bills" }

type ArApRecordModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	TenantID     string    `gorm:"column:tenant_id;index"`
	Type         string    `gorm:"column:type"`
	OrderID      string    `gorm:"column:order_id"`
	Amount       float64   `gorm:"column:amount"`
	Currency     string    `gorm:"column:currency"`
	ExchangeRate float64   `gorm:"column:exchange_rate"`
	AmountCNY    float64   `gorm:"column:amount_cny"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}
func (ArApRecordModel) TableName() string { return "ar_ap_records" }

type CostRecordModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	OrderID   string    `gorm:"column:order_id"`
	SKUID     string    `gorm:"column:sku_id"`
	CostType  string    `gorm:"column:cost_type"`
	Amount    float64   `gorm:"column:amount"`
	Currency  string    `gorm:"column:currency"`
	AmountCNY float64   `gorm:"column:amount_cny"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
func (CostRecordModel) TableName() string { return "cost_records" }

type ProfitReportModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TenantID       string    `gorm:"column:tenant_id;index"`
	OrderID        string    `gorm:"column:order_id"`
	OrderNo        string    `gorm:"column:order_no"`
	SKUID          string    `gorm:"column:sku_id"`
	SKUCode        string    `gorm:"column:sku_code"`
	SaleAmount     float64   `gorm:"column:sale_amount"`
	PurchaseCost   float64   `gorm:"column:purchase_cost"`
	ShippingCost   float64   `gorm:"column:shipping_cost"`
	CommissionCost float64   `gorm:"column:commission_cost"`
	OtherCost      float64   `gorm:"column:other_cost"`
	TotalCost      float64   `gorm:"column:total_cost"`
	GrossProfit    float64   `gorm:"column:gross_profit"`
	ProfitMargin   float64   `gorm:"column:profit_margin"`
	Currency       string    `gorm:"column:currency"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}
func (ProfitReportModel) TableName() string { return "profit_reports" }

type FinanceJournalModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TenantID       string    `gorm:"column:tenant_id;index"`
	OrderID        string    `gorm:"column:order_id"`
	ChangeType     string    `gorm:"column:change_type"`
	Amount         float64   `gorm:"column:amount"`
	BeforeAmount   float64   `gorm:"column:before_amount"`
	AfterAmount    float64   `gorm:"column:after_amount"`
	Currency       string    `gorm:"column:currency"`
	IdempotencyKey string    `gorm:"column:idempotency_key;index"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}
func (FinanceJournalModel) TableName() string { return "finance_journals" }
