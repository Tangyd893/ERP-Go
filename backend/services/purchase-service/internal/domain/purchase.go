package domain

import "time"

// PurchaseStatus 采购单状态
type PurchaseStatus string

const (
	PurchaseDraft     PurchaseStatus = "draft"
	PurchasePending   PurchaseStatus = "pending"
	PurchaseApproved  PurchaseStatus = "approved"
	PurchaseOrdered   PurchaseStatus = "ordered"
	PurchasePartial   PurchaseStatus = "partial"
	PurchaseCompleted PurchaseStatus = "completed"
	PurchaseCancelled PurchaseStatus = "cancelled"
)

// PurchaseOrder 采购单聚合根
type PurchaseOrder struct {
	ID           string          `json:"id"`
	TenantID     string          `json:"tenant_id"`
	SupplierID   string          `json:"supplier_id"`
	SupplierName string          `json:"supplier_name"`
	OrderNo      string          `json:"order_no"`
	Status       PurchaseStatus  `json:"status"`
	Currency     string          `json:"currency"`
	TotalAmount  float64         `json:"total_amount"`
	Items        []*PurchaseItem `json:"items"`
	ExpectedDate time.Time       `json:"expected_date"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// PurchaseItem 采购明细
type PurchaseItem struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	SKUID      string  `json:"sku_id"`
	SKUCode    string  `json:"sku_code"`
	SKUName    string  `json:"sku_name"`
	Quantity   int     `json:"quantity"`
	ReceivedQty int    `json:"received_quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

// Supplier 供应商
type Supplier struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	ContactName string    `json:"contact_name"`
	ContactPhone string   `json:"contact_phone"`
	Email       string    `json:"email"`
	PaymentTerm string    `json:"payment_term"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// InboundOrder 入库单（WMS入库侧）
type InboundOrder struct {
	ID           string         `json:"id"`
	TenantID     string         `json:"tenant_id"`
	PurchaseID   string         `json:"purchase_id"`
	WarehouseID  string         `json:"warehouse_id"`
	Status       string         `json:"status"` // receiving, checking, completed
	Items        []*InboundItem `json:"items"`
	CreatedAt    time.Time      `json:"created_at"`
}

// InboundItem 入库明细
type InboundItem struct {
	ID          string `json:"id"`
	InboundID   string `json:"inbound_id"`
	SKUID       string `json:"sku_id"`
	Quantity    int    `json:"quantity"`
	ReceivedQty int    `json:"received_quantity"`
	PassedQty   int    `json:"passed_quantity"`
	RejectedQty int    `json:"rejected_quantity"`
}
