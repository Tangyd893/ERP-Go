package app

import (
	"context"
	"os"

	"github.com/Tangyd893/ERP-Go/backend/services/report-service/internal/domain"
)

// ReportAppService 报表应用服务（聚合查询型）
type ReportAppService struct {
	demoMode bool
}

func NewReportAppService() *ReportAppService {
	return &ReportAppService{
		demoMode: os.Getenv("DEMO_MODE") == "true" || os.Getenv("ENVIRONMENT") == "",
	}
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
