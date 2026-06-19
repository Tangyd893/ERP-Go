package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/infra/repository"
)

type FinanceAppService struct {
	repo *repository.FinanceRepository
}

func NewFinanceAppService(repo *repository.FinanceRepository) *FinanceAppService {
	return &FinanceAppService{repo: repo}
}

// ── 结算单 ──────────────────────────────────────────────

func (s *FinanceAppService) CreateSettlementBill(ctx context.Context, bill *domain.SettlementBill) error {
	return s.repo.CreateSettlementBill(ctx, bill)
}
func (s *FinanceAppService) ListSettlementBills(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SettlementBill, int64, error) {
	return s.repo.ListSettlementBills(ctx, tenantID, offset, limit)
}

// ImportSettlement 导入平台结算单（含佣金/退款/调整项匹配）
func (s *FinanceAppService) ImportSettlement(ctx context.Context, p domain.SettlementParams) (*domain.SettlementBill, error) {
	net := p.Sales - p.Refunds - p.Commission - p.Fba - p.Other
	bill := &domain.SettlementBill{
		ID: fmt.Sprintf("STL%d", time.Now().UnixNano()), TenantID: p.TenantID,
		StoreID: p.StoreID, PlatformCode: p.Platform, SettlementPeriod: p.Period,
		Currency: p.Currency, TotalSales: p.Sales, TotalRefunds: p.Refunds,
		Commission: p.Commission, FbaFee: p.Fba, OtherFee: p.Other,
		NetAmount: net, Status: "imported", CreatedAt: time.Now(),
	}
	if err := s.repo.CreateSettlementBill(ctx, bill); err != nil {
		return nil, fmt.Errorf("导入结算单失败: %w", err)
	}
	return bill, nil
}

// ── 应收应付 ────────────────────────────────────────────

func (s *FinanceAppService) ListArApRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ArApRecord, int64, error) {
	return s.repo.ListArApRecords(ctx, tenantID, offset, limit)
}

func (s *FinanceAppService) CreateReceivable(ctx context.Context, tenantID, orderID string, amount float64, currency string) (*domain.ArApRecord, error) {
	rate, _ := s.getExchangeRate(ctx, tenantID, currency, "CNY")
	rec := domain.NewReceivable(fmt.Sprintf("AR%d", time.Now().UnixNano()), tenantID, orderID, amount, currency, rate)
	if err := s.repo.CreateArApRecord(ctx, rec); err != nil {
		return nil, fmt.Errorf("创建应收失败: %w", err)
	}
	return rec, nil
}

func (s *FinanceAppService) CreatePayable(ctx context.Context, tenantID, orderID string, amount float64, currency string) (*domain.ArApRecord, error) {
	rate, _ := s.getExchangeRate(ctx, tenantID, currency, "CNY")
	rec := domain.NewPayable(fmt.Sprintf("AP%d", time.Now().UnixNano()), tenantID, orderID, amount, currency, rate)
	if err := s.repo.CreateArApRecord(ctx, rec); err != nil {
		return nil, fmt.Errorf("创建应付失败: %w", err)
	}
	return rec, nil
}

// ── 成本 ────────────────────────────────────────────────

func (s *FinanceAppService) ListCostRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.CostRecord, int64, error) {
	return s.repo.ListCostRecords(ctx, tenantID, offset, limit)
}

// RecordCost 记录成本（采购/物流/佣金等）
func (s *FinanceAppService) RecordCost(ctx context.Context, tenantID, orderID, skuID, costType string, amount float64, currency string) (*domain.CostRecord, error) {
	rate, _ := s.getExchangeRate(ctx, tenantID, currency, "CNY")
	rec := domain.NewCostRecord(fmt.Sprintf("COST%d", time.Now().UnixNano()), domain.CostRecordParams{
		TenantID: tenantID, OrderID: orderID, SKUID: skuID,
		CostType: costType, Amount: amount, Currency: currency, Rate: rate,
	})
	if err := s.repo.CreateCostRecord(ctx, rec); err != nil {
		return nil, fmt.Errorf("记录成本失败: %w", err)
	}
	return rec, nil
}

// ── 利润 ────────────────────────────────────────────────

func (s *FinanceAppService) ListProfitReports(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ProfitReport, int64, error) {
	return s.repo.ListProfitReports(ctx, tenantID, offset, limit)
}

// GenerateProfitReport 生成订单/SKU 利润报表
func (s *FinanceAppService) GenerateProfitReport(ctx context.Context, p domain.ProfitParams) (*domain.ProfitReport, error) {
	rate, _ := s.getExchangeRate(ctx, p.TenantID, p.Currency, "CNY")
	report := &domain.ProfitReport{
		ID: fmt.Sprintf("PR%d", time.Now().UnixNano()), TenantID: p.TenantID,
		OrderID: p.OrderID, OrderNo: p.OrderNo, SKUID: p.SKUID, SKUCode: p.SKUCode,
		SaleAmount: p.SaleAmount, Currency: p.Currency, CreatedAt: time.Now(),
	}
	// 从已记录的成本中汇总
	costs, _, _ := s.repo.ListCostRecords(ctx, p.TenantID, 0, 1000)
	for _, c := range costs {
		if c.OrderID == p.OrderID && (c.SKUID == p.SKUID || p.SKUID == "") {
			switch c.CostType {
			case "purchase": report.PurchaseCost += c.AmountCNY
			case "shipping": report.ShippingCost += c.AmountCNY
			case "commission": report.CommissionCost += c.AmountCNY
			default: report.OtherCost += c.AmountCNY
			}
		}
	}
	// CNY conversion: costs are already in CNY from cost records
	if p.Currency != "CNY" {
		report.PurchaseCost = report.PurchaseCost * rate
		report.ShippingCost = report.ShippingCost * rate
		report.CommissionCost = report.CommissionCost * rate
		report.OtherCost = report.OtherCost * rate
	}
	report.Calculate()
	if err := s.repo.CreateProfitReport(ctx, report); err != nil {
		return nil, fmt.Errorf("生成利润报表失败: %w", err)
	}
	return report, nil
}

// ── 汇率 ────────────────────────────────────────────────

func (s *FinanceAppService) SetExchangeRate(ctx context.Context, tenantID, from, to string, rate float64, source string) (*domain.ExchangeRate, error) {
	r := &domain.ExchangeRate{
		ID: fmt.Sprintf("FX%d", time.Now().UnixNano()), TenantID: tenantID,
		FromCurrency: from, ToCurrency: to, Rate: rate, Source: source, CreatedAt: time.Now(),
	}
	if err := s.repo.CreateExchangeRate(ctx, r); err != nil {
		return nil, fmt.Errorf("设置汇率失败: %w", err)
	}
	return r, nil
}

func (s *FinanceAppService) GetExchangeRate(ctx context.Context, tenantID, from, to string) (*domain.ExchangeRate, error) {
	return s.repo.FindExchangeRate(ctx, tenantID, from, to)
}

func (s *FinanceAppService) getExchangeRate(ctx context.Context, tenantID, from, to string) (float64, error) {
	if from == to { return 1.0, nil }
	r, err := s.repo.FindExchangeRate(ctx, tenantID, from, to)
	if err != nil { return 1.0, nil } // 默认 1:1
	return r.Rate, nil
}

// ── 流水 ────────────────────────────────────────────────

func (s *FinanceAppService) ListJournals(ctx context.Context, tenantID string, offset, limit int) ([]*domain.FinanceJournal, int64, error) {
	return s.repo.ListJournals(ctx, tenantID, offset, limit)
}

func (s *FinanceAppService) RecordJournal(ctx context.Context, p domain.JournalParams) (*domain.FinanceJournal, error) {
	j := domain.NewJournal(fmt.Sprintf("JNL%d", time.Now().UnixNano()), p)
	if err := s.repo.CreateJournal(ctx, j); err != nil {
		return nil, fmt.Errorf("记录流水失败: %w", err)
	}
	return j, nil
}
