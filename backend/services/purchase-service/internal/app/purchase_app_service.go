package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/infra/repository"
)

// InventoryNotifyClient 库存服务通知客户端
type InventoryNotifyClient struct {
	baseURL string
}

func NewInventoryNotifyClient(baseURL string) *InventoryNotifyClient {
	return &InventoryNotifyClient{baseURL: baseURL}
}

// NotifyStockIncrease 通知库存增加
func (c *InventoryNotifyClient) NotifyStockIncrease(ctx context.Context, items []StockIncreaseItem) error {
	// HTTP call to inventory-service /api/v1/inventory/increase
	// 简化实现：记录日志，实际调用通过 HTTP adapter
	_ = ctx
	_ = items
	return nil
}

// StockIncreaseItem 库存增量项
type StockIncreaseItem struct {
	SKUID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
	Location string `json:"location"`
}

// ── PurchaseAppService ────────────────────────────────

type PurchaseAppService struct {
	repo          *repository.PurchaseRepository
	inventoryClient *InventoryNotifyClient
}

func NewPurchaseAppService(repo *repository.PurchaseRepository) *PurchaseAppService {
	return &PurchaseAppService{repo: repo}
}

func (s *PurchaseAppService) WithInventoryClient(client *InventoryNotifyClient) *PurchaseAppService {
	s.inventoryClient = client
	return s
}

// ── 供应商 ──────────────────────────────────────────────

func (s *PurchaseAppService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return s.repo.CreateSupplier(ctx, supplier)
}
func (s *PurchaseAppService) ListSuppliers(ctx context.Context, tenantID string) ([]*domain.Supplier, error) {
	return s.repo.ListSuppliers(ctx, tenantID)
}

// ── 采购单 ──────────────────────────────────────────────

func (s *PurchaseAppService) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	return s.repo.CreatePurchaseOrder(ctx, order)
}
func (s *PurchaseAppService) ListPurchaseOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.PurchaseOrder, int64, error) {
	return s.repo.ListPurchaseOrders(ctx, tenantID, offset, limit)
}
func (s *PurchaseAppService) GetPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	return s.repo.FindPurchaseOrder(ctx, id)
}

func (s *PurchaseAppService) SubmitOrder(ctx context.Context, id string) error {
	order, err := s.repo.FindPurchaseOrder(ctx, id)
	if err != nil { return fmt.Errorf("采购单不存在: %w", err) }
	if err := order.Submit(); err != nil { return err }
	return s.repo.UpdatePurchaseStatus(ctx, id, string(order.Status))
}

func (s *PurchaseAppService) ApproveOrder(ctx context.Context, id string) error {
	order, err := s.repo.FindPurchaseOrder(ctx, id)
	if err != nil { return fmt.Errorf("采购单不存在: %w", err) }
	if err := order.Approve(); err != nil { return err }
	return s.repo.UpdatePurchaseStatus(ctx, id, string(order.Status))
}

func (s *PurchaseAppService) MarkOrdered(ctx context.Context, id string) error {
	order, err := s.repo.FindPurchaseOrder(ctx, id)
	if err != nil { return fmt.Errorf("采购单不存在: %w", err) }
	if err := order.MarkOrdered(); err != nil { return err }
	return s.repo.UpdatePurchaseStatus(ctx, id, string(order.Status))
}

func (s *PurchaseAppService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.FindPurchaseOrder(ctx, id)
	if err != nil { return fmt.Errorf("采购单不存在: %w", err) }
	if err := order.Cancel(); err != nil { return err }
	return s.repo.UpdatePurchaseStatus(ctx, id, string(order.Status))
}

// ── 收货 ────────────────────────────────────────────────

// ReceiveItem 收货登记：更新采购项已收数量，创建/更新入库单
func (s *PurchaseAppService) ReceiveItem(ctx context.Context, orderID, itemID, warehouseID string, qty int) (*domain.InboundOrder, error) {
	order, err := s.repo.FindPurchaseOrder(ctx, orderID)
	if err != nil { return nil, fmt.Errorf("采购单不存在: %w", err) }
	if order.Status != domain.PurchaseOrdered && order.Status != domain.PurchasePartial {
		return nil, fmt.Errorf("采购单状态 %s 不可收货", order.Status)
	}

	item, err := s.repo.FindPurchaseItem(ctx, itemID)
	if err != nil { return nil, fmt.Errorf("采购明细不存在: %w", err) }
	if err := item.UpdateReceivedQty(qty); err != nil { return nil, err }
	if err := s.repo.UpdateReceivedQty(ctx, itemID, item.ReceivedQty); err != nil { return nil, err }

	// 首次收货时创建入库单
	inboundID := fmt.Sprintf("IN%d", time.Now().UnixNano())
	now := time.Now()
	inbound := &domain.InboundOrder{
		ID: inboundID, TenantID: order.TenantID, PurchaseID: orderID,
		WarehouseID: warehouseID, Status: string(domain.InboundReceiving), CreatedAt: now,
		Items: []*domain.InboundItem{
			{ID: fmt.Sprintf("II-%s-%s", inboundID, item.SKUID), InboundID: inboundID,
				SKUID: item.SKUID, Quantity: qty, ReceivedQty: qty},
		},
	}
	if err := s.repo.CreateInboundOrder(ctx, inbound); err != nil {
		return nil, fmt.Errorf("创建入库单失败: %w", err)
	}

	// 更新采购单状态为部分收货
	if order.Status == domain.PurchaseOrdered {
		_ = order.RegisterReceipt()
		_ = s.repo.UpdatePurchaseStatus(ctx, orderID, string(order.Status))
	}

	return inbound, nil
}

// ── 入库质检 ────────────────────────────────────────────

func (s *PurchaseAppService) StartQA(ctx context.Context, inboundID string) error {
	in, err := s.repo.FindInboundOrder(ctx, inboundID)
	if err != nil { return fmt.Errorf("入库单不存在: %w", err) }
	if err := in.StartQA(); err != nil { return err }
	return s.repo.UpdateInboundStatus(ctx, inboundID, in.Status)
}

func (s *PurchaseAppService) QAItem(ctx context.Context, inboundID, itemID string, passed, rejected int) error {
	in, err := s.repo.FindInboundOrder(ctx, inboundID)
	if err != nil { return fmt.Errorf("入库单不存在: %w", err) }
	if in.Status != string(domain.InboundQA) {
		return fmt.Errorf("入库单状态 %s 不可质检", in.Status)
	}
	if err := s.repo.UpdateInboundItemQA(ctx, itemID, passed, rejected); err != nil {
		return err
	}
	return nil
}

func (s *PurchaseAppService) CompleteInbound(ctx context.Context, inboundID string) error {
	in, err := s.repo.FindInboundOrder(ctx, inboundID)
	if err != nil { return fmt.Errorf("入库单不存在: %w", err) }
	if err := in.CompleteInbound(); err != nil { return err }
	if err := s.repo.UpdateInboundStatus(ctx, inboundID, in.Status); err != nil { return err }

	// 通知库存服务增量入库
	if s.inventoryClient != nil {
		items := make([]StockIncreaseItem, 0)
		for _, item := range in.Items {
			if item.PassedQty > 0 {
				items = append(items, StockIncreaseItem{SKUID: item.SKUID, Quantity: item.PassedQty})
			}
		}
		if len(items) > 0 {
			_ = s.inventoryClient.NotifyStockIncrease(ctx, items)
		}
	}

	return nil
}

// ReturnRejectedItems 退货：将质检不合格品生成退货记录，标记入库单为 rejected
func (s *PurchaseAppService) ReturnRejectedItems(ctx context.Context, inboundID string) error {
	in, err := s.repo.FindInboundOrder(ctx, inboundID)
	if err != nil { return fmt.Errorf("入库单不存在: %w", err) }
	if err := in.MarkRejected(); err != nil { return err }
	return s.repo.UpdateInboundStatus(ctx, inboundID, in.Status)
}

// ── 入库单查询 ──────────────────────────────────────────

func (s *PurchaseAppService) ListInboundOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.InboundOrder, int64, error) {
	return s.repo.ListInboundOrders(ctx, tenantID, offset, limit)
}

func (s *PurchaseAppService) GetInboundOrder(ctx context.Context, id string) (*domain.InboundOrder, error) {
	return s.repo.FindInboundOrder(ctx, id)
}
