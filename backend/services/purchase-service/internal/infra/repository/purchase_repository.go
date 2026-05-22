package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/domain"
	"gorm.io/gorm"
)

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
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&models).Error
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
	query := r.db.WithContext(ctx).Model(&PurchaseOrderModel{}).Where("tenant_id = ?", tenantID)
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

func (r *PurchaseRepository) ListInboundOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.InboundOrder, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&InboundOrderModel{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	var models []*InboundOrderModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	orders := make([]*domain.InboundOrder, len(models))
	for i, m := range models {
		orders[i] = &domain.InboundOrder{ID: m.ID, TenantID: m.TenantID, PurchaseID: m.PurchaseID, WarehouseID: m.WarehouseID, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return orders, total, nil
}
