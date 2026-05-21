package domain

import (
	"fmt"
	"time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderPending        OrderStatus = "pending"         // 待审核
	OrderApproved       OrderStatus = "approved"        // 已审核
	OrderLocked         OrderStatus = "locked"          // 已锁定库存
	OrderPicking        OrderStatus = "picking"         // 拣货中
	OrderPacked         OrderStatus = "packed"          // 已打包
	OrderShipped        OrderStatus = "shipped"         // 已发货
	OrderDelivered      OrderStatus = "delivered"       // 已签收
	OrderCompleted      OrderStatus = "completed"       // 已完成
	OrderCancelled      OrderStatus = "cancelled"       // 已取消
	OrderAbnormal       OrderStatus = "abnormal"        // 异常
	OrderPartialShipped OrderStatus = "partial_shipped" // 部分发货
)

// OrderType 订单类型
type OrderType string

const (
	OrderTypeNormal  OrderType = "normal"
	OrderTypeReship  OrderType = "reship"
	OrderTypeReplace OrderType = "replace"
)

// OrderSource 订单来源
type OrderSource string

const (
	OrderSourcePlatform OrderSource = "platform"
	OrderSourceCSV      OrderSource = "csv"
	OrderSourceManual   OrderSource = "manual"
)

// SalesOrder 销售订单聚合根
type SalesOrder struct {
	ID              string       `json:"id"`
	TenantID        string       `json:"tenant_id"`
	StoreID         string       `json:"store_id"`
	PlatformOrderNo string       `json:"platform_order_no"`
	OrderType       OrderType    `json:"order_type"`
	OrderSource     OrderSource  `json:"order_source"`
	Status          OrderStatus  `json:"status"`
	BuyerName       string       `json:"buyer_name"`
	BuyerEmail      string       `json:"buyer_email"`
	Currency        string       `json:"currency"`
	TotalAmount     float64      `json:"total_amount"`
	ShippingFee     float64      `json:"shipping_fee"`
	TaxAmount       float64      `json:"tax_amount"`
	Items           []*OrderItem `json:"items,omitempty"`
	Address         *Address     `json:"address,omitempty"`
	StatusHistory   []*StatusLog `json:"status_history,omitempty"`
	IdempotencyKey  string       `json:"idempotency_key"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

// OrderItem 订单明细
type OrderItem struct {
	ID           string  `json:"id"`
	OrderID      string  `json:"order_id"`
	SKUID        string  `json:"sku_id"`
	SKUCode      string  `json:"sku_code"`
	SKUName      string  `json:"sku_name"`
	PlatformSKU  string  `json:"platform_sku"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
}

// Address 地址值对象
type Address struct {
	ContactName  string `json:"contact_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	State        string `json:"state"`
	City         string `json:"city"`
	District     string `json:"district"`
	StreetLine1  string `json:"street_line1"`
	StreetLine2  string `json:"street_line2"`
	PostalCode   string `json:"postal_code"`
}

// StatusLog 状态变更记录
type StatusLog struct {
	FromStatus OrderStatus `json:"from_status"`
	ToStatus   OrderStatus `json:"to_status"`
	Operator   string      `json:"operator"`
	Remark     string      `json:"remark"`
	CreatedAt  time.Time   `json:"created_at"`
}

// OrderStateMachine 订单状态机
type OrderStateMachine struct {
	transitions map[OrderStatus]map[OrderStatus]bool
}

func NewOrderStateMachine() *OrderStateMachine {
	m := &OrderStateMachine{transitions: make(map[OrderStatus]map[OrderStatus]bool)}
	// 定义合法状态流转
	m.addTransition(OrderPending, OrderApproved)
	m.addTransition(OrderPending, OrderCancelled)
	m.addTransition(OrderPending, OrderAbnormal)

	m.addTransition(OrderApproved, OrderLocked)
	m.addTransition(OrderApproved, OrderAbnormal)
	m.addTransition(OrderApproved, OrderCancelled)

	m.addTransition(OrderLocked, OrderPicking)
	m.addTransition(OrderLocked, OrderAbnormal)
	m.addTransition(OrderLocked, OrderCancelled)

	m.addTransition(OrderPicking, OrderPacked)
	m.addTransition(OrderPicking, OrderAbnormal)

	m.addTransition(OrderPacked, OrderShipped)
	m.addTransition(OrderPacked, OrderAbnormal)

	m.addTransition(OrderShipped, OrderDelivered)
	m.addTransition(OrderShipped, OrderPartialShipped)
	m.addTransition(OrderShipped, OrderAbnormal)

	m.addTransition(OrderPartialShipped, OrderDelivered)
	m.addTransition(OrderPartialShipped, OrderShipped)

	m.addTransition(OrderDelivered, OrderCompleted)

	m.addTransition(OrderAbnormal, OrderPending)

	return m
}

func (m *OrderStateMachine) addTransition(from, to OrderStatus) {
	if m.transitions[from] == nil {
		m.transitions[from] = make(map[OrderStatus]bool)
	}
	m.transitions[from][to] = true
}

// CanTransition 检查状态流转是否合法
func (m *OrderStateMachine) CanTransition(from, to OrderStatus) bool {
	allowed, ok := m.transitions[from]
	if !ok {
		return false
	}
	return allowed[to]
}

// Transition 执行状态流转
func (o *SalesOrder) Transition(target OrderStatus, operator, remark string) error {
	sm := NewOrderStateMachine()
	if !sm.CanTransition(o.Status, target) {
		return fmt.Errorf("订单状态不能从 %s 流转到 %s", o.Status, target)
	}

	o.StatusHistory = append(o.StatusHistory, &StatusLog{
		FromStatus: o.Status,
		ToStatus:   target,
		Operator:   operator,
		Remark:     remark,
		CreatedAt:  time.Now(),
	})
	o.Status = target
	o.UpdatedAt = time.Now()
	return nil
}

// Approve 审核订单
func (o *SalesOrder) Approve(operator string) error {
	return o.Transition(OrderApproved, operator, "审核通过")
}

// Cancel 取消订单，释放库存
func (o *SalesOrder) Cancel(operator, reason string) error {
	return o.Transition(OrderCancelled, operator, reason)
}

// MarkAbnormal 标记异常
func (o *SalesOrder) MarkAbnormal(operator, reason string) error {
	return o.Transition(OrderAbnormal, operator, reason)
}

// CanBeLocked 订单是否可以锁定库存
func (o *SalesOrder) CanBeLocked() bool {
	return o.Status == OrderApproved
}

// CalculateTotal 计算订单总金额
func (o *SalesOrder) CalculateTotal() {
	var total float64
	for _, item := range o.Items {
		item.TotalPrice = float64(item.Quantity) * item.UnitPrice
		total += item.TotalPrice
	}
	o.TotalAmount = total
}

// GetSKUQuantities 获取订单中 SKU 及其数量
func (o *SalesOrder) GetSKUQuantities() map[string]int {
	result := make(map[string]int)
	for _, item := range o.Items {
		result[item.SKUID] += item.Quantity
	}
	return result
}
