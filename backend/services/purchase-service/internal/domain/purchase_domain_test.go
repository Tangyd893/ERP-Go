package domain

import (
	"context"
	"testing"
)

// TestPurchaseToInboundFlow 验证采购单→收货→入库→质检→完成→退货全流程
func TestPurchaseToInboundFlow(t *testing.T) {
	po := createTestPurchaseOrder()
	stepSubmit(t, po)
	stepApprove(t, po)
	stepMarkOrdered(t, po)
	stepReceive(t, po)
	stepCreateInbound(t, po)
	stepReturn(t, po)
	stepComplete(t, po)
}

func createTestPurchaseOrder() *PurchaseOrder {
	return &PurchaseOrder{
		ID: "PO-TEST-001", TenantID: "default", SupplierID: "SP-001",
		OrderNo: "PO-20260101-001", Status: PurchaseDraft,
		Items: []*PurchaseItem{
			{ID: "PI-001", SKUID: "sku-001", SKUCode: "A001", SKUName: "商品A", Quantity: 100, UnitPrice: 9.90, TotalPrice: 990},
		},
	}
}

func stepSubmit(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	check(t, po.Submit(), "提交失败")
	if po.Status != PurchasePending {
		t.Error("应为 pending")
	}
}

func stepApprove(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	check(t, po.Approve(), "审核失败")
	if po.Status != PurchaseApproved {
		t.Error("应为 approved")
	}
}

func stepMarkOrdered(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	check(t, po.MarkOrdered(), "下单失败")
	if po.Status != PurchaseOrdered {
		t.Error("应为 ordered")
	}
}

func stepReceive(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	item := po.Items[0]
	check(t, item.UpdateReceivedQty(50), "收货失败")
	if item.ReceivedQty != 50 {
		t.Errorf("已收应为 50，实际 %d", item.ReceivedQty)
	}
	check(t, po.RegisterReceipt(), "登记收货失败")
	if po.Status != PurchasePartial {
		t.Error("应为 partial")
	}
}

func stepCreateInbound(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	inbound := NewInboundOrder("IN-TEST-001", "default", po.ID, "WH-001")
	inbound.Items = []*InboundItem{
		{ID: "II-001", SKUID: "sku-001", Quantity: 50, ReceivedQty: 50},
	}
	if inbound.Status != string(InboundReceiving) {
		t.Error("应为 receiving")
	}
	check(t, inbound.StartQA(), "开始质检失败")
	if inbound.Status != string(InboundQA) {
		t.Error("应为 qa")
	}
	inbound.Items[0].PassedQty = 48
	inbound.Items[0].RejectedQty = 2
	check(t, inbound.CompleteInbound(), "完成入库失败")
	if inbound.Status != string(InboundPassed) {
		t.Error("应为 passed")
	}
}

func stepReturn(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	inbound2 := NewInboundOrder("IN-TEST-002", "default", po.ID, "WH-001")
	inbound2.Items = []*InboundItem{
		{ID: "II-002", SKUID: "sku-002", Quantity: 10, ReceivedQty: 10, PassedQty: 0, RejectedQty: 10},
	}
	inbound2.StartQA()
	check(t, inbound2.MarkRejected(), "退货失败")
	if inbound2.Status != string(InboundRejected) {
		t.Error("应为 rejected")
	}
}

func stepComplete(t *testing.T, po *PurchaseOrder) {
	t.Helper()
	po.Items[0].ReceivedQty = 100
	check(t, po.Complete(), "完成采购失败")
	if po.Status != PurchaseCompleted {
		t.Error("应为 completed")
	}
}

func check(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
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
