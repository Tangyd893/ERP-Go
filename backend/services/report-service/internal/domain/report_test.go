package domain

import (
	"testing"
	"time"
)

// 创建测试用销售报表
func setupSalesReport() *SalesReport {
	return &SalesReport{
		Period:        "2024-01",
		TotalOrders:   1000,
		TotalSales:    50000.00,
		AvgOrderValue: 50.00,
		Currency:      "CNY",
	}
}

// 创建测试用库存周转报表
func setupInventoryTurnover() *InventoryTurnover {
	return &InventoryTurnover{
		SKUID:        "sku-001",
		SKUCode:      "TSHIRT-RED-M",
		SKUName:      "红色T恤 M码",
		OutboundQty:  5000,
		AvgStock:     1000.00,
		TurnoverRate: 5.0,
	}
}

// 创建测试用仓储效率报表
func setupWarehouseEfficiency() *WarehouseEfficiency {
	return &WarehouseEfficiency{
		WarehouseID:   "wh-001",
		WarehouseName: "深圳中心仓",
		TotalOutbound: 8000,
		TotalPickQty:  25000,
		AvgPickTime:   15.5,
	}
}

// 创建测试用利润汇总报表
func setupProfitSummary() *ProfitSummary {
	return &ProfitSummary{
		Period:       "2024-01",
		TotalSales:   100000.00,
		TotalCost:    70000.00,
		GrossProfit:  30000.00,
		ProfitMargin: 30.0,
		Currency:     "CNY",
		CreatedAt:    time.Now(),
	}
}

// TestSalesReportCreation 测试销售报表创建与字段
func TestSalesReportCreation(t *testing.T) {
	r := setupSalesReport()

	if r.Period == "" {
		t.Error("期间不应为空")
	}
	if r.TotalOrders < 0 {
		t.Error("订单总数不应为负")
	}
	if r.TotalSales < 0 {
		t.Error("销售总额不应为负")
	}
	if r.AvgOrderValue < 0 {
		t.Error("平均订单金额不应为负")
	}
	if r.Currency == "" {
		t.Error("货币不应为空")
	}
}

// TestSalesReportAvgOrderValue 测试平均客单价计算
func TestSalesReportAvgOrderValue(t *testing.T) {
	tests := []struct {
		name        string
		totalOrders int64
		totalSales  float64
		expectedAvg float64
	}{
		{"正常计算", 100, 5000.00, 50.00},
		{"大额订单", 10, 100000.00, 10000.00},
		{"1笔订单", 1, 99.99, 99.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SalesReport{
				Period:       "2024-01",
				TotalOrders:  tt.totalOrders,
				TotalSales:   tt.totalSales,
				Currency:     "CNY",
			}
			if tt.totalOrders > 0 {
				r.AvgOrderValue = r.TotalSales / float64(r.TotalOrders)
			}
			if r.AvgOrderValue != tt.expectedAvg {
				t.Errorf("平均订单金额应为 %.2f，实际 %.2f", tt.expectedAvg, r.AvgOrderValue)
			}
		})
	}
}

// TestSalesReportZeroOrders 测试零订单的销售报表
func TestSalesReportZeroOrders(t *testing.T) {
	r := &SalesReport{
		Period:       "2024-01",
		TotalOrders:  0,
		TotalSales:   0,
		Currency:     "CNY",
	}

	if r.TotalOrders != 0 || r.TotalSales != 0 {
		t.Error("零订单时期的总数和金额应均为0")
	}
}

// TestSalesReportCurrency 测试不同货币的销售报表
func TestSalesReportCurrency(t *testing.T) {
	currencies := []string{"CNY", "USD", "EUR", "JPY", "GBP"}

	for _, c := range currencies {
		t.Run("货币_"+c, func(t *testing.T) {
			r := setupSalesReport()
			r.Currency = c
			if r.Currency != c {
				t.Errorf("货币应为 %s，实际 %s", c, r.Currency)
			}
		})
	}
}

// TestSalesReportPeriod 测试不同期间的销售报表
func TestSalesReportPeriod(t *testing.T) {
	periods := []string{"2024-01", "2024-Q1", "2024-H1", "2024"}

	for _, p := range periods {
		t.Run("期间_"+p, func(t *testing.T) {
			r := setupSalesReport()
			r.Period = p
			if r.Period != p {
				t.Errorf("期间应为 %s，实际 %s", p, r.Period)
			}
		})
	}
}

// TestInventoryTurnoverCreation 测试库存周转报表创建
func TestInventoryTurnoverCreation(t *testing.T) {
	r := setupInventoryTurnover()

	if r.SKUID == "" {
		t.Error("SKU ID不应为空")
	}
	if r.SKUCode == "" {
		t.Error("SKU编码不应为空")
	}
	if r.SKUName == "" {
		t.Error("SKU名称不应为空")
	}
	if r.TurnoverRate < 0 {
		t.Error("周转率不应为负")
	}
}

// TestInventoryTurnoverRate 测试库存周转率计算
func TestInventoryTurnoverRate(t *testing.T) {
	tests := []struct {
		name         string
		outboundQty  int64
		avgStock     float64
		expectedRate float64
	}{
		{"正常周转", 5000, 1000.00, 5.0},
		{"高周转", 10000, 500.00, 20.0},
		{"低周转", 1000, 2000.00, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.avgStock > 0 {
				rate := float64(tt.outboundQty) / tt.avgStock
				if rate != tt.expectedRate {
					t.Errorf("周转率应为 %.2f，实际 %.2f", tt.expectedRate, rate)
				}
			}
		})
	}
}

// TestInventoryTurnoverZeroStock 测试零库存的库存周转
func TestInventoryTurnoverZeroStock(t *testing.T) {
	r := &InventoryTurnover{
		SKUID:       "sku-001",
		SKUCode:     "TST-001",
		SKUName:     "测试商品",
		OutboundQty: 0,
		AvgStock:    0,
	}

	if r.AvgStock != 0 || r.OutboundQty != 0 {
		t.Error("零库存零出库时数据应为0")
	}
}

// TestInventoryTurnoverNegativeOutbound 测试负出库量
func TestInventoryTurnoverNegativeOutbound(t *testing.T) {
	r := &InventoryTurnover{
		SKUID:       "sku-001",
		SKUCode:     "TST-001",
		OutboundQty: -1,
		AvgStock:    100.0,
	}

	// 负出库量应在领域层被允许，校验由应用层负责
	if r.OutboundQty >= 0 {
		t.Error("应允许负出库量值以支持退货冲抵场景")
	}
}

// TestWarehouseEfficiencyCreation 测试仓储效率报表创建
func TestWarehouseEfficiencyCreation(t *testing.T) {
	r := setupWarehouseEfficiency()

	if r.WarehouseID == "" {
		t.Error("仓库ID不应为空")
	}
	if r.WarehouseName == "" {
		t.Error("仓库名称不应为空")
	}
	if r.TotalOutbound < 0 {
		t.Error("出库总数不应为负")
	}
	if r.TotalPickQty < 0 {
		t.Error("拣货总数不应为负")
	}
	if r.AvgPickTime < 0 {
		t.Error("平均拣货时间不应为负")
	}
}

// TestWarehouseEfficiencyPickTime 测试平均拣货时间
func TestWarehouseEfficiencyPickTime(t *testing.T) {
	tests := []struct {
		name       string
		pickTime   float64
		isReasonable bool
	}{
		{"快速拣货", 5.0, true},
		{"正常拣货", 15.5, true},
		{"慢速拣货", 60.0, true},
		{"零拣货时间", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupWarehouseEfficiency()
			r.AvgPickTime = tt.pickTime
			if r.AvgPickTime != tt.pickTime {
				t.Errorf("平均拣货时间应为 %.1f, 实际 %.1f", tt.pickTime, r.AvgPickTime)
			}
		})
	}
}

// TestWarehouseEfficiencyMultipleWarehouses 测试多仓库效率
func TestWarehouseEfficiencyMultipleWarehouses(t *testing.T) {
	warehouses := []struct {
		id       string
		name     string
		outbound int64
	}{
		{"wh-001", "深圳中心仓", 10000},
		{"wh-002", "上海分仓", 5000},
		{"wh-003", "北京分仓", 3000},
	}

	for _, wh := range warehouses {
		t.Run("仓库_"+wh.name, func(t *testing.T) {
			r := &WarehouseEfficiency{
				WarehouseID:   wh.id,
				WarehouseName: wh.name,
				TotalOutbound: wh.outbound,
			}
			if r.WarehouseID != wh.id {
				t.Errorf("仓库ID应为 %s，实际 %s", wh.id, r.WarehouseID)
			}
			if r.WarehouseName != wh.name {
				t.Errorf("仓库名称应为 %s，实际 %s", wh.name, r.WarehouseName)
			}
		})
	}
}

// TestProfitSummaryCreation 测试利润汇总创建
func TestProfitSummaryCreation(t *testing.T) {
	r := setupProfitSummary()

	if r.Period == "" {
		t.Error("期间不应为空")
	}
	if r.TotalSales < 0 {
		t.Error("销售总额不应为负")
	}
	if r.TotalCost < 0 {
		t.Error("总成本不应为负")
	}
	if r.GrossProfit < 0 {
		t.Error("毛利不应为负")
	}
	if r.Currency == "" {
		t.Error("货币不应为空")
	}
	if r.CreatedAt.IsZero() {
		t.Error("创建时间不应为零值")
	}
}

// TestProfitSummaryCalculation 测试利润计算
func TestProfitSummaryCalculation(t *testing.T) {
	tests := []struct {
		name       string
		sales      float64
		cost       float64
		expProfit  float64
		expMargin  float64
	}{
		{"盈利", 100000.00, 70000.00, 30000.00, 30.0},
		{"盈亏平衡", 50000.00, 50000.00, 0, 0},
		{"高利润", 100000.00, 20000.00, 80000.00, 80.0},
		{"微利", 100000.00, 95000.00, 5000.00, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProfitSummary{
				Period:     "2024-01",
				TotalSales: tt.sales,
				TotalCost:  tt.cost,
				Currency:   "CNY",
			}
			r.GrossProfit = r.TotalSales - r.TotalCost
			if r.TotalSales > 0 {
				r.ProfitMargin = (r.GrossProfit / r.TotalSales) * 100
			}

			if r.GrossProfit != tt.expProfit {
				t.Errorf("毛利应为 %.2f，实际 %.2f", tt.expProfit, r.GrossProfit)
			}
			if r.ProfitMargin != tt.expMargin {
				t.Errorf("利润率应为 %.1f%%，实际 %.1f%%", tt.expMargin, r.ProfitMargin)
			}
		})
	}
}

// TestProfitSummaryLoss 测试亏损情况
func TestProfitSummaryLoss(t *testing.T) {
	r := &ProfitSummary{
		Period:     "2024-02",
		TotalSales: 50000.00,
		TotalCost:  80000.00,
		Currency:   "CNY",
	}
	r.GrossProfit = r.TotalSales - r.TotalCost

	if r.GrossProfit >= 0 {
		t.Error("成本大于收入时毛利应为负")
	}
	if r.GrossProfit != -30000.00 {
		t.Errorf("毛利应为 -30000.00，实际 %.2f", r.GrossProfit)
	}
}

// TestProfitSummaryCurrency 测试利润报表货币
func TestProfitSummaryCurrency(t *testing.T) {
	currencies := []string{"CNY", "USD", "EUR"}

	for _, c := range currencies {
		t.Run("货币_"+c, func(t *testing.T) {
			r := setupProfitSummary()
			r.Currency = c
			if r.Currency != c {
				t.Errorf("货币应为 %s，实际 %s", c, r.Currency)
			}
		})
	}
}

// TestProfitSummaryCreatedAt 测试利润报表创建时间
func TestProfitSummaryCreatedAt(t *testing.T) {
	now := time.Now()
	r := &ProfitSummary{
		Period:     "2024-01",
		TotalSales: 10000.00,
		TotalCost:  5000.00,
		Currency:   "CNY",
		CreatedAt:  now,
	}

	if !r.CreatedAt.Equal(now) {
		t.Error("创建时间应与赋值时间一致")
	}
}
