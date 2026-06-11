package domain

import (
	"context"
	"testing"
)

// TestPurchaseToInboundFlow 验证采购单→收货→入库→质检→完成→退货全流程
func TestPurchaseToInboundFlow(t *testing.T) {
	// 1. 创建采购单
	po := &PurchaseOrder{
		ID: "PO-TEST-001", TenantID: "default", SupplierID: "SP-001",
		OrderNo: "PO-20260101-001", Status: PurchaseDraft,
		Items: []*PurchaseItem{
			{ID: "PI-001", SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Quantity: 100, UnitPrice: 9.90, TotalPrice: 990},
		},
	}

	// 2. 提交→审核→下单
	if err := po.Submit(); err != nil {
		t.Fatalf("提交失败: %v", err)
	}
	if po.Status != PurchasePending { t.Error("应为 pending") }

	if err := po.Approve(); err != nil {
		t.Fatalf("审核失败: %v", err)
	}
	if po.Status != PurchaseApproved { t.Error("应为 approved") }

	if err := po.MarkOrdered(); err != nil {
		t.Fatalf("下单失败: %v", err)
	}
	if po.Status != PurchaseOrdered { t.Error("应为 ordered") }

	// 3. 收货
	item := po.Items[0]
	if err := item.UpdateReceivedQty(50); err != nil {
		t.Fatalf("收货失败: %v", err)
	}
	if item.ReceivedQty != 50 { t.Errorf("已收应为 50，实际 %d", item.ReceivedQty) }

	if err := po.RegisterReceipt(); err != nil {
		t.Fatalf("登记收货失败: %v", err)
	}
	if po.Status != PurchasePartial { t.Error("应为 partial") }

	// 4. 创建入库单
	inbound := NewInboundOrder("IN-TEST-001", "default", po.ID, "WH-001")
	inbound.Items = []*InboundItem{
		{ID: "II-001", SKUID: "sku-001", Quantity: 50, ReceivedQty: 50},
	}
	if inbound.Status != string(InboundReceiving) { t.Error("应为 receiving") }

	// 5. 质检
	if err := inbound.StartQA(); err != nil {
		t.Fatalf("开始质检失败: %v", err)
	}
	if inbound.Status != string(InboundQA) { t.Error("应为 qa") }

	inbound.Items[0].PassedQty = 48
	inbound.Items[0].RejectedQty = 2

	// 6. 完成入库
	if err := inbound.CompleteInbound(); err != nil {
		t.Fatalf("完成入库失败: %v", err)
	}
	if inbound.Status != string(InboundPassed) { t.Error("应为 passed") }

	// 7. 退货流程
	inbound2 := NewInboundOrder("IN-TEST-002", "default", po.ID, "WH-001")
	inbound2.Items = []*InboundItem{
		{ID: "II-002", SKUID: "sku-002", Quantity: 10, ReceivedQty: 10, PassedQty: 0, RejectedQty: 10},
	}
	inbound2.StartQA()
	if err := inbound2.MarkRejected(); err != nil {
		t.Fatalf("退货失败: %v", err)
	}
	if inbound2.Status != string(InboundRejected) { t.Error("应为 rejected") }

	// 8. 完成采购
	item.ReceivedQty = 100
	if err := po.Complete(); err != nil {
		t.Fatalf("完成采购失败: %v", err)
	}
	if po.Status != PurchaseCompleted { t.Error("应为 completed") }
}

// TestPurchaseInvalidTransitions 验证非法状态转换
func TestPurchaseInvalidTransitions(t *testing.T) {
	tests := []struct {
		name   string
		status PurchaseStatus
		action func(o *PurchaseOrder) error
	}{
		{"已完成不可取消", PurchaseCompleted, func(o *PurchaseOrder) error { return o.Cancel() }},
		{"已取消不可审核", PurchaseCancelled, func(o *PurchaseOrder) error { return o.Approve() }},
		{"草稿不可直接审核", PurchaseDraft, func(o *PurchaseOrder) error { return o.Approve() }},
		{"已下单不可再提交", PurchaseOrdered, func(o *PurchaseOrder) error { return o.Submit() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			po := &PurchaseOrder{Status: tt.status}
			if err := tt.action(po); err == nil {
				t.Error("应返回错误但未返回")
			}
		})
	}

	// 入库非法转换
	in := &InboundOrder{Status: string(InboundReceiving)}
	if err := in.CompleteInbound(); err == nil {
		t.Error("收货中不可直接完成入库")
	}
}

// TestReceiveExceedQuantity 验证收货不超量
func TestReceiveExceedQuantity(t *testing.T) {
	item := &PurchaseItem{Quantity: 10, ReceivedQty: 0}
	if err := item.UpdateReceivedQty(5); err != nil {
		t.Fatal(err)
	}
	if err := item.UpdateReceivedQty(5); err != nil {
		t.Fatal(err)
	}
	if err := item.UpdateReceivedQty(1); err == nil {
		t.Error("超出订购量应报错")
	}
	if item.ReceivedQty != 10 { t.Errorf("已收应为 10，实际 %d", item.ReceivedQty) }
}

// TestContextCancellation 验证 context 取消
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if ctx.Err() == nil {
		t.Error("已取消的 context 应返回错误")
	}
}
