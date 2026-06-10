package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/domain"
	"gorm.io/gorm"
)

const (
	whereTenantID = "tenant_id = ?"
	orderByDesc   = "created_at DESC"
)

type FinanceRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) CreateSettlementBill(ctx context.Context, bill *domain.SettlementBill) error {
	return r.db.WithContext(ctx).Create(&SettlementBillModel{
		ID: bill.ID, TenantID: bill.TenantID, StoreID: bill.StoreID, PlatformCode: bill.PlatformCode,
		SettlementPeriod: bill.SettlementPeriod, Currency: bill.Currency, TotalSales: bill.TotalSales,
		TotalRefunds: bill.TotalRefunds, Commission: bill.Commission, FbaFee: bill.FbaFee,
		OtherFee: bill.OtherFee, NetAmount: bill.NetAmount, Status: bill.Status, CreatedAt: bill.CreatedAt,
	}).Error
}

func (r *FinanceRepository) ListSettlementBills(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SettlementBill, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&SettlementBillModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*SettlementBillModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	bills := make([]*domain.SettlementBill, len(models))
	for i, m := range models {
		bills[i] = &domain.SettlementBill{
			ID: m.ID, TenantID: m.TenantID, StoreID: m.StoreID, PlatformCode: m.PlatformCode,
			SettlementPeriod: m.SettlementPeriod, Currency: m.Currency, TotalSales: m.TotalSales,
			TotalRefunds: m.TotalRefunds, Commission: m.Commission, FbaFee: m.FbaFee,
			OtherFee: m.OtherFee, NetAmount: m.NetAmount, Status: m.Status, CreatedAt: m.CreatedAt,
		}
	}
	return bills, total, nil
}

func (r *FinanceRepository) ListArApRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ArApRecord, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&ArApRecordModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*ArApRecordModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	records := make([]*domain.ArApRecord, len(models))
	for i, m := range models {
		records[i] = &domain.ArApRecord{ID: m.ID, TenantID: m.TenantID, Type: domain.ArApType(m.Type), OrderID: m.OrderID, Amount: m.Amount, Currency: m.Currency, ExchangeRate: m.ExchangeRate, AmountCNY: m.AmountCNY, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return records, total, nil
}

func (r *FinanceRepository) ListCostRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.CostRecord, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&CostRecordModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*CostRecordModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	records := make([]*domain.CostRecord, len(models))
	for i, m := range models {
		records[i] = &domain.CostRecord{ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, SKUID: m.SKUID, CostType: m.CostType, Amount: m.Amount, Currency: m.Currency, AmountCNY: m.AmountCNY, CreatedAt: m.CreatedAt}
	}
	return records, total, nil
}

func (r *FinanceRepository) ListProfitReports(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ProfitReport, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&ProfitReportModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*ProfitReportModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	reports := make([]*domain.ProfitReport, len(models))
	for i, m := range models {
		reports[i] = &domain.ProfitReport{ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, OrderNo: m.OrderNo, SKUID: m.SKUID, SKUCode: m.SKUCode, SaleAmount: m.SaleAmount, PurchaseCost: m.PurchaseCost, ShippingCost: m.ShippingCost, CommissionCost: m.CommissionCost, OtherCost: m.OtherCost, TotalCost: m.TotalCost, GrossProfit: m.GrossProfit, ProfitMargin: m.ProfitMargin, Currency: m.Currency, CreatedAt: m.CreatedAt}
	}
	return reports, total, nil
}

func (r *FinanceRepository) ListJournals(ctx context.Context, tenantID string, offset, limit int) ([]*domain.FinanceJournal, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&FinanceJournalModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*FinanceJournalModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	journals := make([]*domain.FinanceJournal, len(models))
	for i, m := range models {
		journals[i] = &domain.FinanceJournal{ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, ChangeType: m.ChangeType, Amount: m.Amount, BeforeAmount: m.BeforeAmount, AfterAmount: m.AfterAmount, Currency: m.Currency, IdempotencyKey: m.IdempotencyKey, CreatedAt: m.CreatedAt}
	}
	return journals, total, nil
}
