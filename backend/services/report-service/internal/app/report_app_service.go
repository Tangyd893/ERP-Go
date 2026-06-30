package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/report-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/httpclient"
)

// ReportAppService 报表应用服务（聚合查询型）
type ReportAppService struct {
	demoMode    bool
	orderCli    *httpclient.Client
	productCli  *httpclient.Client
	warehouseCli *httpclient.Client
}

func NewReportAppService() *ReportAppService {
	// T-610: DEMO_MODE 仅在显式设置且非 production 环境时生效
	demoMode := os.Getenv("DEMO_MODE") == "true"
	if demoMode {
		env := os.Getenv("ENVIRONMENT")
		if env == "production" {
			demoMode = false
		}
	}

	orderURL := envOrDefault("ORDER_SERVICE_URL", "http://localhost:8085")
	productURL := envOrDefault("PRODUCT_SERVICE_URL", "http://localhost:8083")
	warehouseURL := envOrDefault("WAREHOUSE_SERVICE_URL", "http://localhost:8087")

	return &ReportAppService{
		demoMode:    demoMode,
		orderCli:    httpclient.New(orderURL),
		productCli:  httpclient.New(productURL),
		warehouseCli: httpclient.New(warehouseURL),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (s *ReportAppService) GetSalesReport(ctx context.Context, tenantID, period string) (*domain.SalesReport, error) {
	if !s.demoMode {
		return s.fetchSalesReport(ctx, tenantID, period)
	}
	return &domain.SalesReport{
		Period: period, TotalOrders: 1280, TotalSales: 56200.50, AvgOrderValue: 43.91, Currency: "USD",
	}, nil
}

func (s *ReportAppService) GetInventoryTurnover(ctx context.Context, tenantID string) ([]*domain.InventoryTurnover, error) {
	if !s.demoMode {
		return s.fetchInventoryTurnover(ctx, tenantID)
	}
	return []*domain.InventoryTurnover{
		{SKUID: "sku-001", SKUCode: "TSHIRT-001", SKUName: "T恤经典款", OutboundQty: 520, AvgStock: 200, TurnoverRate: 2.6},
		{SKUID: "sku-002", SKUCode: "MUG-001", SKUName: "马克杯", OutboundQty: 310, AvgStock: 150, TurnoverRate: 2.07},
	}, nil
}

func (s *ReportAppService) GetWarehouseEfficiency(ctx context.Context, tenantID string) ([]*domain.WarehouseEfficiency, error) {
	if !s.demoMode {
		return s.fetchWarehouseEfficiency(ctx, tenantID)
	}
	return []*domain.WarehouseEfficiency{
		{WarehouseID: "wh-001", WarehouseName: "美东仓", TotalOutbound: 850, TotalPickQty: 3200, AvgPickTime: 12.5},
	}, nil
}

func (s *ReportAppService) GetProfitSummary(ctx context.Context, tenantID, period string) (*domain.ProfitSummary, error) {
	if !s.demoMode {
		return s.fetchProfitSummary(ctx, tenantID, period)
	}
	return &domain.ProfitSummary{
		Period: period, TotalSales: 56200.50, TotalCost: 38500.30,
		GrossProfit: 17700.20, ProfitMargin: 31.5, Currency: "USD",
	}, nil
}

func (s *ReportAppService) fetchSalesReport(ctx context.Context, tenantID, period string) (*domain.SalesReport, error) {
	return &domain.SalesReport{
		Period: period, TotalOrders: 0, TotalSales: 0, AvgOrderValue: 0, Currency: "USD",
	}, nil
}

func (s *ReportAppService) fetchInventoryTurnover(ctx context.Context, tenantID string) ([]*domain.InventoryTurnover, error) {
	return []*domain.InventoryTurnover{}, nil
}

func (s *ReportAppService) fetchWarehouseEfficiency(ctx context.Context, tenantID string) ([]*domain.WarehouseEfficiency, error) {
	return []*domain.WarehouseEfficiency{}, nil
}

func (s *ReportAppService) fetchProfitSummary(ctx context.Context, tenantID, period string) (*domain.ProfitSummary, error) {
	return &domain.ProfitSummary{
		Period: period, TotalSales: 0, TotalCost: 0, GrossProfit: 0, ProfitMargin: 0, Currency: "USD",
	}, nil
}

// GetDashboard 看板聚合数据（T-636）
func (s *ReportAppService) GetDashboard(ctx context.Context, tenantID string) (*domain.DashboardData, error) {
	if !s.demoMode {
		return s.fetchDashboard(ctx, tenantID)
	}
	return &domain.DashboardData{
		OrderCount:    328,
		SalesAmount:   12850.00,
		OutboundCount: 186,
		SkuCount:      1256,
		Trend: []*domain.TrendPoint{
			{Date: "周一", OrderCount: 45, SalesAmount: 2100},
			{Date: "周二", OrderCount: 52, SalesAmount: 3200},
			{Date: "周三", OrderCount: 38, SalesAmount: 1850},
			{Date: "周四", OrderCount: 65, SalesAmount: 4200},
			{Date: "周五", OrderCount: 48, SalesAmount: 2800},
			{Date: "周六", OrderCount: 32, SalesAmount: 1500},
			{Date: "周日", OrderCount: 28, SalesAmount: 1200},
		},
		Timeliness: domain.TimelinessRate{
			Within24h: 152,
			Within48h: 28,
			Overdue:   6,
		},
	}, nil
}

// OrderItem 订单列表中的单项（用于提取 total_amount）
type OrderItem struct {
	TotalAmount float64 `json:"total_amount"`
	CreatedAt   string  `json:"created_at"`
}

// OutboundItem 出库单列表中的单项（用于提取时效信息）
type OutboundItem struct {
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *ReportAppService) fetchDashboard(ctx context.Context, tenantID string) (*domain.DashboardData, error) {
	d := &domain.DashboardData{}

	// 注入 tenant 到 HTTP 客户端
	orderCli := s.orderCli.WithTenant(tenantID)
	productCli := s.productCli.WithTenant(tenantID)
	warehouseCli := s.warehouseCli.WithTenant(tenantID)

	// 1. 并发获取各服务计数（使用 page_size=1 仅取 total）
	type countResult struct {
		name  string
		total int64
		err   error
	}
	results := make(chan countResult, 3)

	go func() {
		total, err := orderCli.GetPageTotal(ctx, "/api/v1/orders?page=1&page_size=1")
		results <- countResult{"order", total, err}
	}()
	go func() {
		total, err := productCli.GetPageTotal(ctx, "/api/v1/skus?page=1&page_size=1")
		results <- countResult{"product", total, err}
	}()
	go func() {
		total, err := warehouseCli.GetPageTotal(ctx, "/api/v1/outbounds?page=1&page_size=1")
		results <- countResult{"warehouse", total, err}
	}()

	for i := 0; i < 3; i++ {
		r := <-results
		if r.err != nil {
			return nil, fmt.Errorf("report: 获取 %s 计数失败: %w", r.name, r.err)
		}
		switch r.name {
		case "order":
			d.OrderCount = r.total
		case "product":
			d.SkuCount = r.total
		case "warehouse":
			d.OutboundCount = r.total
		}
	}

	// 2. 获取订单详情（用于销售额聚合 + 趋势）
	orders, err := fetchOrders(ctx, orderCli)
	if err != nil {
		// 订单数据失败不阻塞整体返回，KPI 已有 count
		orders = nil
	}

	if len(orders) > 0 {
		var totalSales float64
		trendMap := make(map[string]*domain.TrendPoint)
		for _, o := range orders {
			totalSales += o.TotalAmount
			date := extractDate(o.CreatedAt)
			if date == "" {
				continue
			}
			if tp, ok := trendMap[date]; ok {
				tp.OrderCount++
				tp.SalesAmount += o.TotalAmount
			} else {
				trendMap[date] = &domain.TrendPoint{
					Date: date, OrderCount: 1, SalesAmount: o.TotalAmount,
				}
			}
		}
		d.SalesAmount = totalSales

		// 趋势：取最近 7 天，按日期排序
		dates := sortedKeys(trendMap)
		if len(dates) > 7 {
			dates = dates[len(dates)-7:]
		}
		d.Trend = make([]*domain.TrendPoint, 0, len(dates))
		for _, date := range dates {
			d.Trend = append(d.Trend, trendMap[date])
		}
	}

	// 3. 获取出库单详情（用于出库及时率）
	outbounds, err := fetchOutbounds(ctx, warehouseCli)
	if err != nil {
		outbounds = nil
	}

	if len(outbounds) > 0 {
		now := time.Now()
		var within24, within48, overdue int64
		for _, ob := range outbounds {
			created := parseTime(ob.CreatedAt)
			if created.IsZero() {
				continue
			}
			hours := now.Sub(created).Hours()
			switch {
			case hours <= 24:
				within24++
			case hours <= 48:
				within48++
			default:
				overdue++
			}
		}
		// 未出库的也算 overdue
		for _, ob := range outbounds {
			if ob.Status != "shipped" && ob.Status != "cancelled" {
				updated := parseTime(ob.UpdatedAt)
				if !updated.IsZero() {
					hours := time.Since(updated).Hours()
					if hours > 48 {
						overdue++
					} else if hours > 24 {
						within48++
					} else {
						within24++
					}
				}
			}
		}
		d.Timeliness = domain.TimelinessRate{
			Within24h: within24,
			Within48h: within48,
			Overdue:   overdue,
		}
	}

	return d, nil
}

func fetchOrders(ctx context.Context, cli *httpclient.Client) ([]OrderItem, error) {
	list, _, err := cli.GetList(ctx, "/api/v1/orders?page=1&page_size=200")
	if err != nil {
		return nil, err
	}
	var orders []OrderItem
	if err := json.Unmarshal(list, &orders); err != nil {
		return nil, fmt.Errorf("report: 解析订单列表失败: %w", err)
	}
	return orders, nil
}

func fetchOutbounds(ctx context.Context, cli *httpclient.Client) ([]OutboundItem, error) {
	list, _, err := cli.GetList(ctx, "/api/v1/outbounds?page=1&page_size=200")
	if err != nil {
		return nil, err
	}
	var outbounds []OutboundItem
	if err := json.Unmarshal(list, &outbounds); err != nil {
		return nil, fmt.Errorf("report: 解析出库单列表失败: %w", err)
	}
	return outbounds, nil
}

func extractDate(ts string) string {
	if len(ts) >= 10 {
		return ts[:10] // "2025-01-15T..." → "2025-01-15"
	}
	return ""
}

func parseTime(ts string) time.Time {
	t, _ := time.Parse(time.RFC3339, ts)
	return t
}

func sortedKeys(m map[string]*domain.TrendPoint) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
