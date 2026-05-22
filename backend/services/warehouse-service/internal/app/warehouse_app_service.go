package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/infra/repository"
)

type WarehouseAppService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseAppService(repo *repository.WarehouseRepository) *WarehouseAppService {
	return &WarehouseAppService{repo: repo}
}

func (s *WarehouseAppService) CreateOutbound(ctx context.Context, order *domain.OutboundOrder) error {
	return s.repo.CreateOutbound(ctx, order)
}
func (s *WarehouseAppService) ListOutbounds(ctx context.Context, tenantID string, offset, limit int) ([]*domain.OutboundOrder, int64, error) {
	return s.repo.ListOutbounds(ctx, tenantID, offset, limit)
}
func (s *WarehouseAppService) GetOutbound(ctx context.Context, id string) (*domain.OutboundOrder, error) {
	return s.repo.FindOutbound(ctx, id)
}
func (s *WarehouseAppService) UpdateOutboundStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateOutboundStatus(ctx, id, status)
}
func (s *WarehouseAppService) ListPickTasks(ctx context.Context, outboundID string) ([]*domain.PickTask, error) {
	return s.repo.ListPickTasks(ctx, outboundID)
}
func (s *WarehouseAppService) PickScan(ctx context.Context, taskID string, pickedQty int) error {
	return s.repo.UpdatePickQty(ctx, taskID, pickedQty, "picked")
}
func (s *WarehouseAppService) ListWarehouses(ctx context.Context, tenantID string) ([]*domain.Warehouse, error) {
	return s.repo.ListWarehouses(ctx, tenantID)
}
