package domain

import (
	"fmt"
	"time"
)

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

// ── 采购单业务规则 ──────────────────────────────────────

// Submit 提交审核
func (o *PurchaseOrder) Submit() error {
	if o.Status != PurchaseDraft {
		return fmt.Errorf("采购单状态 %s 不可提交", o.Status)
	}
	o.Status = PurchasePending
	o.UpdatedAt = time.Now()
	return nil
}

// Approve 审核通过
func (o *PurchaseOrder) Approve() error {
	if o.Status != PurchasePending {
		return fmt.Errorf("采购单状态 %s 不可审核", o.Status)
	}
	o.Status = PurchaseApproved
	o.UpdatedAt = time.Now()
	return nil
}

// MarkOrdered 确认下单给供应商
func (o *PurchaseOrder) MarkOrdered() error {
	if o.Status != PurchaseApproved {
		return fmt.Errorf("采购单状态 %s 不可下单", o.Status)
	}
	o.Status = PurchaseOrdered
	o.UpdatedAt = time.Now()
	return nil
}

// RegisterReceipt 登记收货（部分或全部）
func (o *PurchaseOrder) RegisterReceipt() error {
	if o.Status != PurchaseOrdered && o.Status != PurchasePartial {
		return fmt.Errorf("采购单状态 %s 不可收货", o.Status)
	}
	o.Status = PurchasePartial
	o.UpdatedAt = time.Now()
	return nil
}

// Complete 完成采购
func (o *PurchaseOrder) Complete() error {
	if o.Status != PurchasePartial {
		return fmt.Errorf("采购单状态 %s 不可完成", o.Status)
	}
	o.Status = PurchaseCompleted
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel 取消采购
func (o *PurchaseOrder) Cancel() error {
	if o.Status == PurchaseCompleted || o.Status == PurchaseCancelled {
		return fmt.Errorf("采购单状态 %s 不可取消", o.Status)
	}
	o.Status = PurchaseCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// UpdateReceivedQty 更新采购明细已收数量
func (item *PurchaseItem) UpdateReceivedQty(qty int) error {
	if qty <= 0 {
		return fmt.Errorf("收货数量必须大于 0")
	}
	if item.ReceivedQty+qty > item.Quantity {
		return fmt.Errorf("收货数量 %d 超过订购量 %d", item.ReceivedQty+qty, item.Quantity)
	}
	item.ReceivedQty += qty
	return nil
}

// ── 入库单业务规则 ──────────────────────────────────────

// InboundStatus 入库单状态
type InboundStatus string

const (
	InboundReceiving InboundStatus = "receiving"  // 收货中
	InboundQA        InboundStatus = "qa"         // 质检中
	InboundPassed    InboundStatus = "passed"     // 合格入库
	InboundRejected  InboundStatus = "rejected"   // 退货
)

// NewInboundOrder 创建入库单
func NewInboundOrder(id, tenantID, purchaseID, warehouseID string) *InboundOrder {
	return &InboundOrder{
		ID:          id,
		TenantID:    tenantID,
		PurchaseID:  purchaseID,
		WarehouseID: warehouseID,
		Status:      string(InboundReceiving),
		CreatedAt:   time.Now(),
	}
}

// StartQA 开始质检
func (in *InboundOrder) StartQA() error {
	if in.Status != string(InboundReceiving) {
		return fmt.Errorf("入库单状态 %s 不可开始质检", in.Status)
	}
	in.Status = string(InboundQA)
	return nil
}

// MarkRejected 标记退货
func (in *InboundOrder) MarkRejected() error {
	if in.Status != string(InboundQA) {
		return fmt.Errorf("入库单状态 %s 不可退货", in.Status)
	}
	in.Status = string(InboundRejected)
	return nil
}

// CompleteInbound 完成入库（合格入库）
func (in *InboundOrder) CompleteInbound() error {
	if in.Status != string(InboundQA) {
		return fmt.Errorf("入库单状态 %s 不可完成入库", in.Status)
	}
	in.Status = string(InboundPassed)
	return nil
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
