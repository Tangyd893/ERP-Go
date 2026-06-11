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
	task, err := s.repo.FindPickTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("拣货任务不存在: %w", err)
	}
	// 幂等：已是 picked/checked 状态直接返回成功
	if task.Status == "picked" || task.Status == "checked" || task.Status == "check_pending" {
		return nil
	}
	if task.Status != "pending" && task.Status != "picking" {
		return fmt.Errorf("拣货任务状态 %s 不可拣货", task.Status)
	}
	if pickedQty <= 0 || pickedQty > task.Quantity {
		return fmt.Errorf("拣货数量 %d 无效（应 > 0 且 ≤ %d）", pickedQty, task.Quantity)
	}
	if err := s.repo.UpdatePickQty(ctx, taskID, pickedQty, "picked"); err != nil {
		return err
	}
	return nil
}

// CheckScan 复核扫码：校验 SKU + 数量匹配
func (s *WarehouseAppService) CheckScan(ctx context.Context, outboundID, skuID string, qty int) error {
	// 幂等：检查出库单状态
	ob, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return fmt.Errorf("出库单不存在: %w", err)
	}
	if ob.Status == domain.OutboundChecked || ob.Status == domain.OutboundPacking ||
		ob.Status == domain.OutboundPacked || ob.Status == domain.OutboundWeighed ||
		ob.Status == domain.OutboundShipped {
		return nil // 已复核或更后状态，幂等返回成功
	}
	if ob.Status != domain.OutboundPicked && ob.Status != domain.OutboundChecking {
		return fmt.Errorf("出库单状态 %s 不可复核", ob.Status)
	}

	// 校验 SKU 属于此出库单且数量匹配
	item, err := s.repo.FindOutboundItem(ctx, outboundID, skuID)
	if err != nil {
		return fmt.Errorf("出库明细不存在（SKU: %s）: %w", skuID, err)
	}
	if qty <= 0 || qty > item.Quantity {
		return fmt.Errorf("复核数量 %d 无效（应 > 0 且 ≤ %d）", qty, item.Quantity)
	}

	// 更新已复核数量
	if err := s.repo.UpdateCheckedQty(ctx, item.ID, qty); err != nil {
		return err
	}
	// 状态推进到复核中
	if ob.Status == domain.OutboundPicked {
		_ = s.repo.UpdateOutboundStatus(ctx, outboundID, string(domain.OutboundChecking))
	}
	return nil
}

// CompleteCheck 完成复核（所有 SKU 复核完毕后调用）
func (s *WarehouseAppService) CompleteCheck(ctx context.Context, outboundID string) error {
	ob, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return err
	}
	if err := ob.CompleteChecking(); err != nil {
		return err
	}
	return s.repo.UpdateOutboundStatus(ctx, outboundID, string(ob.Status))
}

// Pack 打包并持久化包裹记录
func (s *WarehouseAppService) Pack(ctx context.Context, outboundID string, weight float64) (*domain.Package, error) {
	ob, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return nil, fmt.Errorf("出库单不存在: %w", err)
	}
	// 幂等：已打包或更后状态
	if ob.Status == domain.OutboundPacked || ob.Status == domain.OutboundWeighed ||
		ob.Status == domain.OutboundShipped {
		return &domain.Package{OutboundID: outboundID, Weight: weight}, nil
	}
	if err := ob.CompletePacking(); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateOutboundStatus(ctx, outboundID, string(ob.Status)); err != nil {
		return nil, err
	}
	now := time.Now()
	pkg := &domain.Package{
		ID:         fmt.Sprintf("PKG%d", now.UnixNano()),
		OutboundID: outboundID,
		Weight:     weight,
		CreatedAt:  now,
	}
	if err := s.repo.CreatePackage(ctx, pkg); err != nil {
		return nil, fmt.Errorf("创建包裹记录失败: %w", err)
	}
	return pkg, nil
}

// Weigh 称重并记录重量
func (s *WarehouseAppService) Weigh(ctx context.Context, outboundID string, weight float64) error {
	ob, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return fmt.Errorf("出库单不存在: %w", err)
	}
	// 幂等：已称重或已发货
	if ob.Status == domain.OutboundWeighed || ob.Status == domain.OutboundShipped {
		return nil
	}
	if err := ob.Weigh(); err != nil {
		return err
	}
	if err := s.repo.UpdateOutboundStatus(ctx, outboundID, string(ob.Status)); err != nil {
		return err
	}
	// 更新包裹重量
	now := time.Now()
	pkg := &domain.Package{
		ID:         fmt.Sprintf("PKG%d", now.UnixNano()),
		OutboundID: outboundID,
		Weight:     weight,
		CreatedAt:  now,
	}
	_ = s.repo.CreatePackage(ctx, pkg) // 称重时若打包阶段未创建包裹则补建
	return nil
}

// ── Wave 波次管理 ──────────────────────────────────────

// CreateWave 创建拣货波次
func (s *WarehouseAppService) CreateWave(ctx context.Context, warehouseID, name string, outboundIDs []string) (*domain.Wave, error) {
	// 校验每个出库单存在且状态为 created
	for _, obID := range outboundIDs {
		ob, err := s.repo.FindOutbound(ctx, obID)
		if err != nil {
			return nil, fmt.Errorf("出库单 %s 不存在: %w", obID, err)
		}
		if ob.Status != domain.OutboundCreated && ob.Status != domain.OutboundWaved {
			return nil, fmt.Errorf("出库单 %s 状态为 %s，不可加入波次", obID, ob.Status)
		}
	}

	wave := domain.NewWave(fmt.Sprintf("WAVE%d", time.Now().UnixNano()), warehouseID, name)
	for _, obID := range outboundIDs {
		if err := wave.AddOutbound(obID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.CreateWave(ctx, wave); err != nil {
		return nil, fmt.Errorf("创建波次失败: %w", err)
	}
	return wave, nil
}

// GetWave 查询波次详情
func (s *WarehouseAppService) GetWave(ctx context.Context, id string) (*domain.Wave, error) {
	return s.repo.FindWave(ctx, id)
}

// ListWaves 列出仓库波次
func (s *WarehouseAppService) ListWaves(ctx context.Context, warehouseID string) ([]*domain.Wave, error) {
	return s.repo.ListWaves(ctx, warehouseID)
}

// StartWave 开始波次拣货
func (s *WarehouseAppService) StartWave(ctx context.Context, id string) error {
	wave, err := s.repo.FindWave(ctx, id)
	if err != nil {
		return fmt.Errorf("波次不存在: %w", err)
	}
	if err := wave.StartPicking(); err != nil {
		return err
	}
	// 更新所有关联出库单状态为 picking
	for _, obID := range wave.OutboundIDs {
		ob, err := s.repo.FindOutbound(ctx, obID)
		if err != nil {
			continue
		}
		if err := ob.StartPicking(); err != nil {
			continue
		}
		_ = s.repo.UpdateOutboundStatus(ctx, obID, string(ob.Status))
	}
	return s.repo.UpdateWaveStatus(ctx, id, wave.Status)
}

// MarkOutboundAbnormal 出库异常上报
func (s *WarehouseAppService) MarkOutboundAbnormal(ctx context.Context, outboundID, reason string) error {
	ob, err := s.repo.FindOutbound(ctx, outboundID)
	if err != nil {
		return fmt.Errorf("出库单不存在: %w", err)
	}
	if err := ob.MarkAbnormal(reason); err != nil {
		return err
	}
	return s.repo.UpdateOutboundStatus(ctx, outboundID, string(ob.Status))
}

// CompleteWave 完成波次
func (s *WarehouseAppService) CompleteWave(ctx context.Context, id string) error {
	wave, err := s.repo.FindWave(ctx, id)
	if err != nil {
		return fmt.Errorf("波次不存在: %w", err)
	}
	if err := wave.Complete(); err != nil {
		return err
	}
	return s.repo.UpdateWaveStatus(ctx, id, wave.Status)
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
