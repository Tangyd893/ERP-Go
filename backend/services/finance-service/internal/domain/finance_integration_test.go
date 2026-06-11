package domain

import (
	"math"
	"testing"
)

// TestFinanceFullFlow 验证结算→应收→成本→利润→流水全流程
func TestFinanceFullFlow(t *testing.T) {
	// 1. 汇率
	rate := 7.25

	// 2. 创建应收（平台销售）
	ar := NewReceivable("AR-001", "default", "order-001", 100.0, "USD", rate)
	if ar.Type != ArApReceivable {
		t.Error("应为 receivable")
	}
	if ar.AmountCNY != 725.0 {
		t.Errorf("CNY 应为 725.0，实际 %.2f", ar.AmountCNY)
	}

	// 3. 创建应付（采购成本）
	ap := NewPayable("AP-001", "default", "order-001", 50.0, "USD", rate)
	if ap.Type != ArApPayable {
		t.Error("应为 payable")
	}
	if ap.AmountCNY != 362.5 {
		t.Errorf("CNY 应为 362.5，实际 %.2f", ap.AmountCNY)
	}

	// 4. 记录成本
	purchase := NewCostRecord("C-001", "default", "order-001", "sku-001", "purchase", 30.0, "USD", rate)
	shipping := NewCostRecord("C-002", "default", "order-001", "sku-001", "shipping", 5.0, "USD", rate)
	commission := NewCostRecord("C-003", "default", "order-001", "sku-001", "commission", 8.0, "USD", rate)

	if purchase.AmountCNY != 217.5 {
		t.Errorf("采购成本 CNY 应为 217.5，实际 %.2f", purchase.AmountCNY)
	}

	// 5. 生成利润报表
	report := &ProfitReport{
		ID: "PR-001", TenantID: "default", OrderID: "order-001", OrderNo: "SO-001",
		SKUID: "sku-001", SKUCode: "A001", SaleAmount: 100.0, Currency: "USD",
		PurchaseCost: purchase.AmountCNY, ShippingCost: shipping.AmountCNY,
		CommissionCost: commission.AmountCNY, OtherCost: 0,
	}
	report.Calculate()

	// purchase: 30*7.25=217.5, shipping: 5*7.25=36.25, commission: 8*7.25=58.0
	expectedCost := 217.5 + 36.25 + 58.0
	if report.TotalCost != expectedCost {
		t.Errorf("总成本应为 %.2f，实际 %.2f", expectedCost, report.TotalCost)
	}
	expectedProfit := 100.0 - expectedCost
	if report.GrossProfit != expectedProfit {
		t.Errorf("毛利应为 %.2f，实际 %.2f", expectedProfit, report.GrossProfit)
	}
	if report.ProfitMargin <= 0 {
		t.Log("利润率为负（成本 > 售价），符合预期")
	}

	// 6. 财务流水
	j := NewJournal("J-001", "default", "order-001", "cost_record", 36.25, 0, 253.75, "USD", "idem-001")
	if j.ChangeType != "cost_record" {
		t.Error("流水类型应为 cost_record")
	}
	if j.BeforeAmount != 0 {
		t.Error("变动前金额应为 0")
	}
	if j.AfterAmount != 253.75 {
		t.Errorf("变动后金额应为 253.75，实际 %.2f", j.AfterAmount)
	}
}

// TestProfitReportPositive 验证正向利润
func TestProfitReportPositive(t *testing.T) {
	report := &ProfitReport{
		SaleAmount: 100.0, PurchaseCost: 30.0, ShippingCost: 5.0,
		CommissionCost: 10.0, OtherCost: 0, Currency: "CNY",
	}
	report.Calculate()
	if report.TotalCost != 45.0 {
		t.Errorf("总成本应为 45.0，实际 %.2f", report.TotalCost)
	}
	if report.GrossProfit != 55.0 {
		t.Errorf("毛利应为 55.0，实际 %.2f", report.GrossProfit)
	}
	if math.Abs(report.ProfitMargin-55.0) > 0.001 {
		t.Errorf("利润率应为 55%%，实际 %.2f%%", report.ProfitMargin)
	}
}

// TestExchangeRateDefaults 验证汇率默认 1:1
func TestExchangeRateDefaults(t *testing.T) {
	// 同币种
	ar := NewReceivable("AR-CNY", "default", "o-001", 100.0, "CNY", 1.0)
	if ar.AmountCNY != 100.0 {
		t.Errorf("CNY→CNY 应为 100.0，实际 %.2f", ar.AmountCNY)
	}
}
