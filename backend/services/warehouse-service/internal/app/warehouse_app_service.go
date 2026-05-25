package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/infra/repository"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
)

type WarehouseAppService struct {
	repo              *repository.WarehouseRepository
	fulfillmentClient *OrderFulfillmentClient
}

func NewWarehouseAppService(repo *repository.WarehouseRepository) *WarehouseAppService {
	return &WarehouseAppService{repo: repo}
}

func (s *WarehouseAppService) WithFulfillmentClient(client *OrderFulfillmentClient) *WarehouseAppService {
	s.fulfillmentClient = client
	return s
}

func (s *WarehouseAppService) CreateOutbound(ctx context.Context, order *domain.OutboundOrder) error {
	if len(order.Items) == 0 {
		return s.repo.CreateOutbound(ctx, order)
	}
	now := time.Now()
	for i, item := range order.Items {
		if item.ID == "" {
			item.ID = fmt.Sprintf("OI%d-%d", now.UnixNano(), i)
		}
	}
	if order.Status == "" {
		order.Status = domain.OutboundPicking
	}
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
	if err := s.repo.UpdatePickQty(ctx, taskID, pickedQty, "picked"); err != nil {
		return err
	}
	return nil
}

// ConfirmShip 出库确认：更新状态并回调 Order 履约
func (s *WarehouseAppService) ConfirmShip(ctx context.Context, outboundID, trackingNo, carrier string) error {
	outbound, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateOutboundStatus(ctx, outboundID, string(domain.OutboundShipped)); err != nil {
		return err
	}
	if s.fulfillmentClient == nil {
		return nil
	}
	items := make([]workflows.OrderItemData, 0, len(outbound.Items))
	for _, it := range outbound.Items {
		qty := it.Quantity
		if it.PickedQty > 0 {
			qty = it.PickedQty
		}
		items = append(items, workflows.OrderItemData{
			SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Qty: qty,
		})
	}
	return s.fulfillmentClient.NotifyOutboundShipped(ctx, workflows.OutboundShippedData{
		OutboundID:  outbound.ID,
		OrderID:     outbound.OrderID,
		TenantID:    outbound.TenantID,
		WarehouseID: outbound.WarehouseID,
		Items:       items,
		TrackingNo:  trackingNo,
		Carrier:     carrier,
	})
}
func (s *WarehouseAppService) ListWarehouses(ctx context.Context, tenantID string) ([]*domain.Warehouse, error) {
	return s.repo.ListWarehouses(ctx, tenantID)
}
