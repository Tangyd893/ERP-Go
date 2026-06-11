package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/domain"
	"gorm.io/gorm"
)

const whereTenantID = "tenant_id = ?"

type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreateSupplier(ctx context.Context, s *domain.Supplier) error {
	return r.db.WithContext(ctx).Create(&SupplierModel{
		ID: s.ID, TenantID: s.TenantID, Name: s.Name, Code: s.Code,
		ContactName: s.ContactName, ContactPhone: s.ContactPhone, Email: s.Email,
		PaymentTerm: s.PaymentTerm, Status: s.Status, CreatedAt: s.CreatedAt,
	}).Error
}

func (r *PurchaseRepository) ListSuppliers(ctx context.Context, tenantID string) ([]*domain.Supplier, error) {
	var models []*SupplierModel
	err := r.db.WithContext(ctx).Where(whereTenantID, tenantID).Find(&models).Error
	if err != nil { return nil, err }
	suppliers := make([]*domain.Supplier, len(models))
	for i, m := range models {
		suppliers[i] = &domain.Supplier{ID: m.ID, TenantID: m.TenantID, Name: m.Name, Code: m.Code, ContactName: m.ContactName, ContactPhone: m.ContactPhone, Email: m.Email, PaymentTerm: m.PaymentTerm, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return suppliers, nil
}

func (r *PurchaseRepository) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		om := &PurchaseOrderModel{
			ID: order.ID, TenantID: order.TenantID, SupplierID: order.SupplierID, SupplierName: order.SupplierName,
			OrderNo: order.OrderNo, Status: string(order.Status), Currency: order.Currency, TotalAmount: order.TotalAmount,
			ExpectedDate: order.ExpectedDate, CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
		}
		if err := tx.Create(om).Error; err != nil { return err }
		for _, item := range order.Items {
			if err := tx.Create(&PurchaseItemModel{
				ID: item.ID, OrderID: order.ID, SKUID: item.SKUID, SKUCode: item.SKUCode, SKUName: item.SKUName,
				Quantity: item.Quantity, ReceivedQty: item.ReceivedQty, UnitPrice: item.UnitPrice, TotalPrice: item.TotalPrice,
			}).Error; err != nil { return err }
		}
		return nil
	})
}

func (r *PurchaseRepository) ListPurchaseOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.PurchaseOrder, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&PurchaseOrderModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*PurchaseOrderModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	orders := make([]*domain.PurchaseOrder, len(models))
	for i, m := range models {
		var items []*PurchaseItemModel
		r.db.WithContext(ctx).Where("order_id = ?", m.ID).Find(&items)
		domainItems := make([]*domain.PurchaseItem, len(items))
		for j, it := range items {
			domainItems[j] = &domain.PurchaseItem{ID: it.ID, OrderID: m.ID, SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Quantity: it.Quantity, ReceivedQty: it.ReceivedQty, UnitPrice: it.UnitPrice, TotalPrice: it.TotalPrice}
		}
		orders[i] = &domain.PurchaseOrder{ID: m.ID, TenantID: m.TenantID, SupplierID: m.SupplierID, SupplierName: m.SupplierName, OrderNo: m.OrderNo, Status: domain.PurchaseStatus(m.Status), Currency: m.Currency, TotalAmount: m.TotalAmount, Items: domainItems, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
	}
	return orders, total, nil
}

func (r *PurchaseRepository) UpdatePurchaseStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&PurchaseOrderModel{}).Where("id = ?", id).Update("status", status).Error
}

// FindPurchaseOrder 按 ID 查询采购单
func (r *PurchaseRepository) FindPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	var m PurchaseOrderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	var items []*PurchaseItemModel
	r.db.WithContext(ctx).Where("order_id = ?", m.ID).Find(&items)
	domainItems := make([]*domain.PurchaseItem, len(items))
	for i, it := range items {
		domainItems[i] = &domain.PurchaseItem{ID: it.ID, OrderID: m.ID, SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Quantity: it.Quantity, ReceivedQty: it.ReceivedQty, UnitPrice: it.UnitPrice, TotalPrice: it.TotalPrice}
	}
	return &domain.PurchaseOrder{ID: m.ID, TenantID: m.TenantID, SupplierID: m.SupplierID, SupplierName: m.SupplierName, OrderNo: m.OrderNo, Status: domain.PurchaseStatus(m.Status), Currency: m.Currency, TotalAmount: m.TotalAmount, Items: domainItems, ExpectedDate: m.ExpectedDate, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

// FindPurchaseItem 按 ID 查询采购明细
func (r *PurchaseRepository) FindPurchaseItem(ctx context.Context, itemID string) (*domain.PurchaseItem, error) {
	var m PurchaseItemModel
	if err := r.db.WithContext(ctx).Where("id = ?", itemID).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.PurchaseItem{ID: m.ID, OrderID: m.OrderID, SKUID: m.SKUID, SKUCode: m.SKUCode, SKUName: m.SKUName, Quantity: m.Quantity, ReceivedQty: m.ReceivedQty, UnitPrice: m.UnitPrice, TotalPrice: m.TotalPrice}, nil
}

// UpdateReceivedQty 更新采购项已收数量
func (r *PurchaseRepository) UpdateReceivedQty(ctx context.Context, itemID string, qty int) error {
	return r.db.WithContext(ctx).Model(&PurchaseItemModel{}).Where("id = ?", itemID).Update("received_quantity", qty).Error
}

// CreateInboundOrder 创建入库单（含明细）
func (r *PurchaseRepository) CreateInboundOrder(ctx context.Context, in *domain.InboundOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&InboundOrderModel{
			ID: in.ID, TenantID: in.TenantID, PurchaseID: in.PurchaseID,
			WarehouseID: in.WarehouseID, Status: in.Status, CreatedAt: in.CreatedAt,
		}).Error; err != nil {
			return err
		}
		for _, item := range in.Items {
			if err := tx.Create(&InboundItemModel{
				ID: item.ID, InboundID: in.ID, SKUID: item.SKUID,
				Quantity: item.Quantity, ReceivedQty: item.ReceivedQty,
				PassedQty: item.PassedQty, RejectedQty: item.RejectedQty,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindInboundOrder 查询入库单
func (r *PurchaseRepository) FindInboundOrder(ctx context.Context, id string) (*domain.InboundOrder, error) {
	var m InboundOrderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	var items []*InboundItemModel
	r.db.WithContext(ctx).Where("inbound_id = ?", m.ID).Find(&items)
	domainItems := make([]*domain.InboundItem, len(items))
	for i, it := range items {
		domainItems[i] = &domain.InboundItem{ID: it.ID, InboundID: m.ID, SKUID: it.SKUID, Quantity: it.Quantity, ReceivedQty: it.ReceivedQty, PassedQty: it.PassedQty, RejectedQty: it.RejectedQty}
	}
	return &domain.InboundOrder{ID: m.ID, TenantID: m.TenantID, PurchaseID: m.PurchaseID, WarehouseID: m.WarehouseID, Status: m.Status, Items: domainItems, CreatedAt: m.CreatedAt}, nil
}

// UpdateInboundStatus 更新入库单状态
func (r *PurchaseRepository) UpdateInboundStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&InboundOrderModel{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateInboundItemQA 更新入库明细质检结果
func (r *PurchaseRepository) UpdateInboundItemQA(ctx context.Context, itemID string, passed, rejected int) error {
	return r.db.WithContext(ctx).Model(&InboundItemModel{}).Where("id = ?", itemID).Updates(map[string]interface{}{
		"passed_quantity":   passed,
		"rejected_quantity": rejected,
	}).Error
}

func (r *PurchaseRepository) ListInboundOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.InboundOrder, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&InboundOrderModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*InboundOrderModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	orders := make([]*domain.InboundOrder, len(models))
	for i, m := range models {
		orders[i] = &domain.InboundOrder{ID: m.ID, TenantID: m.TenantID, PurchaseID: m.PurchaseID, WarehouseID: m.WarehouseID, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return orders, total, nil
}
