package domain

import (
	"math"
	"testing"
)

func TestProfitReportCalculateTotalCost(t *testing.T) {
	p := &ProfitReport{
		PurchaseCost:   10.0,
		ShippingCost:   5.0,
		CommissionCost: 2.5,
		OtherCost:      1.0,
	}
	p.CalculateTotalCost()

	expected := 18.5
	if math.Abs(p.TotalCost-expected) > 0.001 {
		t.Errorf("总成本应为 %.2f，实际: %.2f", expected, p.TotalCost)
	}
}

func TestProfitReportCalculateProfit(t *testing.T) {
	p := &ProfitReport{
		SaleAmount:   30.0,
		PurchaseCost: 10.0,
		ShippingCost: 5.0,
		TotalCost:    15.0,
	}
	p.CalculateProfit()

	expectedProfit := 15.0
	if math.Abs(p.GrossProfit-expectedProfit) > 0.001 {
		t.Errorf("毛利润应为 %.2f，实际: %.2f", expectedProfit, p.GrossProfit)
	}

	expectedMargin := 50.0
	if math.Abs(p.ProfitMargin-expectedMargin) > 0.001 {
		t.Errorf("利润率应为 %.2f%%，实际: %.2f%%", expectedMargin, p.ProfitMargin)
	}
}

func TestProfitReportZeroSaleAmount(t *testing.T) {
	p := &ProfitReport{
		SaleAmount:   0,
		PurchaseCost: 10.0,
		TotalCost:    10.0,
	}
	p.CalculateProfit()

	if p.GrossProfit != -10.0 {
		t.Errorf("销售额为0时利润应为 -10.0，实际: %.2f", p.GrossProfit)
	}
	if p.ProfitMargin != 0 {
		t.Errorf("销售额为0时利润率应为 0，实际: %.2f", p.ProfitMargin)
	}
}

func TestProfitReportCalculate(t *testing.T) {
	p := &ProfitReport{
		SaleAmount:     100.0,
		PurchaseCost:   40.0,
		ShippingCost:   15.0,
		CommissionCost: 5.0,
		OtherCost:      3.0,
	}
	p.Calculate()

	if p.TotalCost != 63.0 {
		t.Errorf("总成本应为 63.0，实际: %.2f", p.TotalCost)
	}
	if p.GrossProfit != 37.0 {
		t.Errorf("毛利润应为 37.0，实际: %.2f", p.GrossProfit)
	}
	if p.ProfitMargin != 37.0 {
		t.Errorf("利润率应为 37.0%%，实际: %.2f%%", p.ProfitMargin)
	}
}
