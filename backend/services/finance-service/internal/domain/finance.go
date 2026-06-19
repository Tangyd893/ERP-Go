package domain

import "time"

// ── 参数结构体（降低函数签名复杂度）─────────────────────

// SettlementParams 结算单导入参数
type SettlementParams struct {
	TenantID, StoreID, Platform, Period, Currency string
	Sales, Refunds, Commission, Fba, Other        float64
}

// ProfitParams 利润报表生成参数
type ProfitParams struct {
	TenantID, OrderID, OrderNo, SKUID, SKUCode, Currency string
	SaleAmount                                           float64
}

// CostRecordParams 成本记录参数
type CostRecordParams struct {
	TenantID, OrderID, SKUID, CostType, Currency string
	Amount, Rate                                 float64
}

// JournalParams 财务流水参数
type JournalParams struct {
	TenantID, OrderID, ChangeType, Currency, IdempotencyKey string
	Amount, Before, After                                   float64
}

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

// ExchangeRate 汇率
type ExchangeRate struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency string    `json:"to_currency"`
	Rate       float64   `json:"rate"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
}

// ── 业务方法 ──────────────────────────────────────────────

// NewReceivable 创建应收记录
func NewReceivable(id, tenantID, orderID string, amount float64, currency string, rate float64) *ArApRecord {
	return &ArApRecord{
		ID: id, TenantID: tenantID, Type: ArApReceivable, OrderID: orderID,
		Amount: amount, Currency: currency, ExchangeRate: rate,
		AmountCNY: amount * rate, Status: "pending", CreatedAt: time.Now(),
	}
}

// NewPayable 创建应付记录
func NewPayable(id, tenantID, orderID string, amount float64, currency string, rate float64) *ArApRecord {
	return &ArApRecord{
		ID: id, TenantID: tenantID, Type: ArApPayable, OrderID: orderID,
		Amount: amount, Currency: currency, ExchangeRate: rate,
		AmountCNY: amount * rate, Status: "pending", CreatedAt: time.Now(),
	}
}

// NewCostRecord 创建成本记录
func NewCostRecord(id string, p CostRecordParams) *CostRecord {
	return &CostRecord{
		ID: id, TenantID: p.TenantID, OrderID: p.OrderID, SKUID: p.SKUID,
		CostType: p.CostType, Amount: p.Amount, Currency: p.Currency,
		AmountCNY: p.Amount * p.Rate, CreatedAt: time.Now(),
	}
}

// NewJournal 创建财务流水
func NewJournal(id string, p JournalParams) *FinanceJournal {
	return &FinanceJournal{
		ID: id, TenantID: p.TenantID, OrderID: p.OrderID, ChangeType: p.ChangeType,
		Amount: p.Amount, BeforeAmount: p.Before, AfterAmount: p.After,
		Currency: p.Currency, IdempotencyKey: p.IdempotencyKey, CreatedAt: time.Now(),
	}
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
