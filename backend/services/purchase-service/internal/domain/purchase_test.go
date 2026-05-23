package domain

import (
	"testing"
	"time"
)

// 创建测试用采购单
func setupPurchaseOrder() *PurchaseOrder {
	return &PurchaseOrder{
		ID:           "po-001",
		TenantID:     "tenant-001",
		SupplierID:   "sup-001",
		SupplierName: "优质供应商有限公司",
		OrderNo:      "PO-2024-0001",
		Status:       PurchaseDraft,
		Currency:     "CNY",
		ExpectedDate: time.Now().Add(7 * 24 * time.Hour),
		Items: []*PurchaseItem{
			{ID: "item-1", SKUID: "sku-001", SKUCode: "TSHIRT-RED-M", SKUName: "红色T恤 M码", Quantity: 100, UnitPrice: 25.00, TotalPrice: 2500.00},
			{ID: "item-2", SKUID: "sku-002", SKUCode: "MUG-WHITE", SKUName: "白色马克杯", Quantity: 200, UnitPrice: 8.00, TotalPrice: 1600.00},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// 创建测试用供应商
func setupSupplier() *Supplier {
	return &Supplier{
		ID:           "sup-001",
		TenantID:     "tenant-001",
		Name:         "优质供应商有限公司",
		Code:         "SUP001",
		ContactName:  "张三",
		ContactPhone: "13800138000",
		Email:        "zhangsan@supplier.com",
		PaymentTerm:  "30天",
		Status:       "active",
		CreatedAt:    time.Now(),
	}
}

// 创建测试用入库单
func setupInboundOrder() *InboundOrder {
	return &InboundOrder{
		ID:          "inb-001",
		TenantID:    "tenant-001",
		PurchaseID:  "po-001",
		WarehouseID: "wh-001",
		Status:      "receiving",
		Items: []*InboundItem{
			{ID: "inbitem-1", SKUID: "sku-001", Quantity: 100, ReceivedQty: 0, PassedQty: 0, RejectedQty: 0},
			{ID: "inbitem-2", SKUID: "sku-002", Quantity: 200, ReceivedQty: 0, PassedQty: 0, RejectedQty: 0},
		},
		CreatedAt: time.Now(),
	}
}

// TestPurchaseStatus 测试采购单状态常量
func TestPurchaseStatus(t *testing.T) {
	tests := []struct {
		name   string
		status PurchaseStatus
		want   string
	}{
		{"草稿", PurchaseDraft, "draft"},
		{"待审核", PurchasePending, "pending"},
		{"已审核", PurchaseApproved, "approved"},
		{"已下单", PurchaseOrdered, "ordered"},
		{"部分收货", PurchasePartial, "partial"},
		{"已完成", PurchaseCompleted, "completed"},
		{"已取消", PurchaseCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("状态值应为 %s，实际 %s", tt.want, tt.status)
			}
		})
	}
}

// TestPurchaseOrderCreation 测试采购单创建
func TestPurchaseOrderCreation(t *testing.T) {
	po := setupPurchaseOrder()

	if po.ID == "" {
		t.Error("采购单ID不应为空")
	}
	if po.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if po.SupplierID == "" {
		t.Error("供应商ID不应为空")
	}
	if po.SupplierName == "" {
		t.Error("供应商名称不应为空")
	}
	if po.OrderNo == "" {
		t.Error("采购单号不应为空")
	}
	if po.Status != PurchaseDraft {
		t.Errorf("新建采购单状态应为 draft，实际 %s", po.Status)
	}
}

// TestPurchaseOrderItems 测试采购单明细
func TestPurchaseOrderItems(t *testing.T) {
	po := setupPurchaseOrder()

	if len(po.Items) != 2 {
		t.Errorf("采购明细应为2项，实际 %d", len(po.Items))
	}
}

// TestPurchaseOrderTotalAmount 测试采购单总金额
func TestPurchaseOrderTotalAmount(t *testing.T) {
	po := setupPurchaseOrder()
	po.TotalAmount = 0
	for _, item := range po.Items {
		po.TotalAmount += item.TotalPrice
	}

	expected := 2500.00 + 1600.00 // 4100.00
	if po.TotalAmount != expected {
		t.Errorf("采购总金额应为 %.2f，实际 %.2f", expected, po.TotalAmount)
	}
}

// TestPurchaseOrderEmptyItems 测试空明细的采购单
func TestPurchaseOrderEmptyItems(t *testing.T) {
	po := &PurchaseOrder{
		ID:           "po-empty",
		TenantID:     "tenant-001",
		SupplierID:   "sup-001",
		SupplierName: "测试供应商",
		OrderNo:      "PO-2024-EMPTY",
		Status:       PurchaseDraft,
		Items:        []*PurchaseItem{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if len(po.Items) != 0 {
		t.Errorf("空明细采购单 items 应为0，实际 %d", len(po.Items))
	}
}

// TestPurchaseItemCreation 测试采购明细创建
func TestPurchaseItemCreation(t *testing.T) {
	item := &PurchaseItem{
		ID:          "item-001",
		OrderID:     "po-001",
		SKUID:       "sku-001",
		SKUCode:     "TSHIRT-RED-M",
		SKUName:     "红色T恤 M码",
		Quantity:    50,
		ReceivedQty: 0,
		UnitPrice:   25.00,
		TotalPrice:  1250.00,
	}

	if item.SKUID == "" {
		t.Error("SKU ID不应为空")
	}
	if item.Quantity <= 0 {
		t.Error("数量应大于0")
	}
	if item.UnitPrice <= 0 {
		t.Error("单价应大于0")
	}
}

// TestPurchaseItemTotalPrice 测试采购明细金额计算
func TestPurchaseItemTotalPrice(t *testing.T) {
	item := &PurchaseItem{
		ID:       "item-001",
		Quantity: 100,
		UnitPrice: 25.00,
	}
	item.TotalPrice = float64(item.Quantity) * item.UnitPrice

	expected := 2500.00
	if item.TotalPrice != expected {
		t.Errorf("明细金额应为 %.2f，实际 %.2f", expected, item.TotalPrice)
	}
}

// TestPurchaseItemReceivedQty 测试已收货数量
func TestPurchaseItemReceivedQty(t *testing.T) {
	item := &PurchaseItem{
		ID:          "item-001",
		Quantity:    100,
		ReceivedQty: 80,
		UnitPrice:   25.00,
	}

	if item.ReceivedQty > item.Quantity {
		t.Error("已收货数量不应超过采购数量")
	}
}

// TestPurchaseStatusTransitions 测试采购单状态流转
func TestPurchaseStatusTransitions(t *testing.T) {
	transitions := []struct {
		from PurchaseStatus
		to   PurchaseStatus
	}{
		{PurchaseDraft, PurchasePending},
		{PurchasePending, PurchaseApproved},
		{PurchaseApproved, PurchaseOrdered},
		{PurchaseOrdered, PurchasePartial},
		{PurchasePartial, PurchaseCompleted},
		{PurchaseDraft, PurchaseCancelled},
		{PurchasePending, PurchaseCancelled},
	}

	for _, tr := range transitions {
		po := &PurchaseOrder{
			ID:     "test",
			Status: tr.from,
		}
		po.Status = tr.to
		if po.Status != tr.to {
			t.Errorf("状态从 %s 流转到 %s 后应为 %s", tr.from, tr.to, tr.to)
		}
	}
}

// TestSupplierCreation 测试供应商创建
func TestSupplierCreation(t *testing.T) {
	supplier := setupSupplier()

	if supplier.ID == "" {
		t.Error("供应商ID不应为空")
	}
	if supplier.Name == "" {
		t.Error("供应商名称不应为空")
	}
	if supplier.Code == "" {
		t.Error("供应商编码不应为空")
	}
	if supplier.ContactName == "" {
		t.Error("联系人不应为空")
	}
	if supplier.ContactPhone == "" {
		t.Error("联系电话不应为空")
	}
	if supplier.Status != "active" {
		t.Errorf("状态应为 active，实际 %s", supplier.Status)
	}
}

// TestSupplierPaymentTerm 测试供应商付款条件
func TestSupplierPaymentTerm(t *testing.T) {
	tests := []struct {
		name string
		term string
	}{
		{"月结30天", "30天"},
		{"月结60天", "60天"},
		{"现结", "现结"},
		{"预付", "预付"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupSupplier()
			s.PaymentTerm = tt.term
			if s.PaymentTerm != tt.term {
				t.Errorf("付款条件应为 %s，实际 %s", tt.term, s.PaymentTerm)
			}
		})
	}
}

// TestSupplierStatus 测试供应商状态
func TestSupplierStatus(t *testing.T) {
	statuses := []string{"active", "inactive", "blacklisted"}

	s := setupSupplier()
	for _, status := range statuses {
		s.Status = status
		if s.Status != status {
			t.Errorf("状态应为 %s，实际 %s", status, s.Status)
		}
	}
}

// TestSupplierWithoutContact 测试无联系方式的供应商
func TestSupplierWithoutContact(t *testing.T) {
	s := &Supplier{
		ID:       "sup-no-contact",
		TenantID: "tenant-001",
		Name:     "无联系方式供应商",
		Code:     "SUP002",
		Status:   "active",
	}

	if s.ContactName != "" {
		t.Error("无联系方式时联系人应为空")
	}
	if s.Email != "" {
		t.Error("无联系方式时邮箱应为空")
	}
}

// TestInboundOrderCreation 测试入库单创建
func TestInboundOrderCreation(t *testing.T) {
	inbound := setupInboundOrder()

	if inbound.ID == "" {
		t.Error("入库单ID不应为空")
	}
	if inbound.PurchaseID == "" {
		t.Error("关联采购单ID不应为空")
	}
	if inbound.WarehouseID == "" {
		t.Error("仓库ID不应为空")
	}
	if inbound.Status != "receiving" {
		t.Errorf("入库单初始状态应为 receiving，实际 %s", inbound.Status)
	}
}

// TestInboundOrderStatus 测试入库单状态
func TestInboundOrderStatus(t *testing.T) {
	statuses := []string{"receiving", "checking", "completed"}

	inbound := setupInboundOrder()
	for _, s := range statuses {
		inbound.Status = s
		if inbound.Status != s {
			t.Errorf("入库状态应为 %s，实际 %s", s, inbound.Status)
		}
	}
}

// TestInboundItemCreation 测试入库明细创建
func TestInboundItemCreation(t *testing.T) {
	item := &InboundItem{
		ID:          "inbitem-001",
		InboundID:   "inb-001",
		SKUID:       "sku-001",
		Quantity:    100,
		ReceivedQty: 95,
		PassedQty:   90,
		RejectedQty: 5,
	}

	if item.SKUID == "" {
		t.Error("入库明细 SKU ID不应为空")
	}
	if item.Quantity <= 0 {
		t.Error("入库数量应大于0")
	}
	if item.PassedQty+item.RejectedQty > item.ReceivedQty {
		t.Error("合格+不合格数量不应超过收货数量")
	}
	if item.ReceivedQty > item.Quantity {
		t.Error("收货数量不应超过应入库数量")
	}
}

// TestInboundItemReceivedValidation 测试入库明细收货校验
func TestInboundItemReceivedValidation(t *testing.T) {
	tests := []struct {
		name        string
		quantity    int
		receivedQty int
		passedQty   int
		rejectedQty int
		valid       bool
	}{
		{"正常收货", 100, 90, 85, 5, true},
		{"全部合格", 100, 100, 100, 0, true},
		{"全检不合格", 100, 100, 0, 100, true},
		{"收货超量", 100, 120, 100, 20, false},
		{"质检超收货", 100, 80, 70, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &InboundItem{
				ID:          "test-item",
				SKUID:       "sku-001",
				Quantity:    tt.quantity,
				ReceivedQty: tt.receivedQty,
				PassedQty:   tt.passedQty,
				RejectedQty: tt.rejectedQty,
			}

			valid := item.ReceivedQty <= item.Quantity &&
				item.PassedQty+item.RejectedQty <= item.ReceivedQty

			if valid != tt.valid {
				t.Errorf("校验结果应为 %v，实际 %v", tt.valid, valid)
			}
		})
	}
}

// TestPurchaseOrderCurrency 测试采购单货币
func TestPurchaseOrderCurrency(t *testing.T) {
	currencies := []string{"CNY", "USD", "EUR"}

	po := setupPurchaseOrder()
	for _, c := range currencies {
		po.Currency = c
		if po.Currency != c {
			t.Errorf("货币应为 %s，实际 %s", c, po.Currency)
		}
	}
}

// TestInboundOrderEmptyItems 测试无明细入库单
func TestInboundOrderEmptyItems(t *testing.T) {
	inbound := &InboundOrder{
		ID:          "inb-empty",
		TenantID:    "tenant-001",
		PurchaseID:  "po-001",
		WarehouseID: "wh-001",
		Status:      "receiving",
		Items:       []*InboundItem{},
		CreatedAt:   time.Now(),
	}

	if len(inbound.Items) != 0 {
		t.Errorf("空入库单明细应为0，实际 %d", len(inbound.Items))
	}
}
