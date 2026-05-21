package domain

import "time"

// 应收应付类型
type ArApType string

const (
	ArApReceivable ArApType = "receivable" // 应收
	ArApPayable    ArApType = "payable"    // 应付
)

// SettlementBill 结算单聚合根
type SettlementBill struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	StoreID        string    `json:"store_id"`
	PlatformCode   string    `json:"platform_code"`
	SettlementPeriod string  `json:"settlement_period"`
	Currency       string    `json:"currency"`
	TotalSales     float64   `json:"total_sales"`
	TotalRefunds   float64   `json:"total_refunds"`
	Commission     float64   `json:"commission"`
	FbaFee         float64   `json:"fba_fee"`
	OtherFee       float64   `json:"other_fee"`
	NetAmount      float64   `json:"net_amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// ArApRecord 应收应付记录
type ArApRecord struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Type        ArApType  `json:"type"`
	OrderID     string    `json:"order_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	ExchangeRate float64  `json:"exchange_rate"`
	AmountCNY   float64   `json:"amount_cny"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// CostRecord 成本记录
type CostRecord struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	OrderID    string    `json:"order_id"`
	SKUID      string    `json:"sku_id"`
	CostType   string    `json:"cost_type"` // purchase, shipping, commission, fba, other
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	AmountCNY  float64   `json:"amount_cny"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProfitReport 利润报表
type ProfitReport struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	OrderID        string    `json:"order_id"`
	OrderNo        string    `json:"order_no"`
	SKUID          string    `json:"sku_id"`
	SKUCode        string    `json:"sku_code"`
	SaleAmount     float64   `json:"sale_amount"`
	PurchaseCost   float64   `json:"purchase_cost"`
	ShippingCost   float64   `json:"shipping_cost"`
	CommissionCost float64   `json:"commission_cost"`
	OtherCost      float64   `json:"other_cost"`
	TotalCost      float64   `json:"total_cost"`
	GrossProfit    float64   `json:"gross_profit"`
	ProfitMargin   float64   `json:"profit_margin"`
	Currency       string    `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
}

// FinanceJournal 财务流水
type FinanceJournal struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	OrderID       string    `json:"order_id"`
	ChangeType    string    `json:"change_type"`
	Amount        float64   `json:"amount"`
	BeforeAmount  float64   `json:"before_amount"`
	AfterAmount   float64   `json:"after_amount"`
	Currency      string    `json:"currency"`
	IdempotencyKey string   `json:"idempotency_key"`
	CreatedAt     time.Time `json:"created_at"`
}
