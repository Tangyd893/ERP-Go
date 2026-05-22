package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/infra/repository"
)

type PurchaseAppService struct {
	repo *repository.PurchaseRepository
}

func NewPurchaseAppService(repo *repository.PurchaseRepository) *PurchaseAppService {
	return &PurchaseAppService{repo: repo}
}

func (s *PurchaseAppService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return s.repo.CreateSupplier(ctx, supplier)
}
func (s *PurchaseAppService) ListSuppliers(ctx context.Context, tenantID string) ([]*domain.Supplier, error) {
	return s.repo.ListSuppliers(ctx, tenantID)
}
func (s *PurchaseAppService) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	return s.repo.CreatePurchaseOrder(ctx, order)
}
func (s *PurchaseAppService) ListPurchaseOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.PurchaseOrder, int64, error) {
	return s.repo.ListPurchaseOrders(ctx, tenantID, offset, limit)
}
func (s *PurchaseAppService) ListInboundOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.InboundOrder, int64, error) {
	return s.repo.ListInboundOrders(ctx, tenantID, offset, limit)
}
