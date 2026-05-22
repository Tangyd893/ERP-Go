package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/infra/repository"
)

type FinanceAppService struct {
	repo *repository.FinanceRepository
}

func NewFinanceAppService(repo *repository.FinanceRepository) *FinanceAppService {
	return &FinanceAppService{repo: repo}
}

func (s *FinanceAppService) CreateSettlementBill(ctx context.Context, bill *domain.SettlementBill) error {
	return s.repo.CreateSettlementBill(ctx, bill)
}
func (s *FinanceAppService) ListSettlementBills(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SettlementBill, int64, error) {
	return s.repo.ListSettlementBills(ctx, tenantID, offset, limit)
}
func (s *FinanceAppService) ListArApRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ArApRecord, int64, error) {
	return s.repo.ListArApRecords(ctx, tenantID, offset, limit)
}
func (s *FinanceAppService) ListCostRecords(ctx context.Context, tenantID string, offset, limit int) ([]*domain.CostRecord, int64, error) {
	return s.repo.ListCostRecords(ctx, tenantID, offset, limit)
}
func (s *FinanceAppService) ListProfitReports(ctx context.Context, tenantID string, offset, limit int) ([]*domain.ProfitReport, int64, error) {
	return s.repo.ListProfitReports(ctx, tenantID, offset, limit)
}
func (s *FinanceAppService) ListJournals(ctx context.Context, tenantID string, offset, limit int) ([]*domain.FinanceJournal, int64, error) {
	return s.repo.ListJournals(ctx, tenantID, offset, limit)
}
