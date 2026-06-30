package domain

import "time"

// SalesReport 销售报表
type SalesReport struct {
	Period       string  `json:"period"`
	TotalOrders  int64   `json:"total_orders"`
	TotalSales   float64 `json:"total_sales"`
	AvgOrderValue float64 `json:"avg_order_value"`
	Currency     string  `json:"currency"`
}

// InventoryTurnover 库存周转报表
type InventoryTurnover struct {
	SKUID       string  `json:"sku_id"`
	SKUCode     string  `json:"sku_code"`
	SKUName     string  `json:"sku_name"`
	OutboundQty int64   `json:"outbound_qty"`
	AvgStock    float64 `json:"avg_stock"`
	TurnoverRate float64 `json:"turnover_rate"`
}

// WarehouseEfficiency 仓储效率报表
type WarehouseEfficiency struct {
	WarehouseID   string  `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	TotalOutbound int64   `json:"total_outbound"`
	TotalPickQty  int64   `json:"total_pick_qty"`
	AvgPickTime   float64 `json:"avg_pick_time_minutes"`
}

// ProfitSummary 利润汇总报表
type ProfitSummary struct {
	Period        string  `json:"period"`
	TotalSales    float64 `json:"total_sales"`
	TotalCost     float64 `json:"total_cost"`
	GrossProfit   float64 `json:"gross_profit"`
	ProfitMargin  float64 `json:"profit_margin"`
	Currency      string  `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date        string  `json:"date"`
	OrderCount  int64   `json:"order_count"`
	SalesAmount float64 `json:"sales_amount"`
}

// TimelinessRate 出库及时率
type TimelinessRate struct {
	Within24h int64 `json:"within_24h"`
	Within48h int64 `json:"within_48h"`
	Overdue   int64 `json:"overdue"`
}

// DashboardData 看板聚合数据（T-636）
type DashboardData struct {
	OrderCount    int64          `json:"order_count"`
	SalesAmount   float64        `json:"sales_amount"`
	OutboundCount int64          `json:"outbound_count"`
	SkuCount      int64          `json:"sku_count"`
	Trend         []*TrendPoint  `json:"trend"`
	Timeliness    TimelinessRate `json:"timeliness"`
}
