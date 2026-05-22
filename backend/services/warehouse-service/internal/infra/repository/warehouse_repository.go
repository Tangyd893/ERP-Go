package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/domain"
	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) CreateOutbound(ctx context.Context, order *domain.OutboundOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		om := &OutboundOrderModel{
			ID: order.ID, TenantID: order.TenantID, OrderID: order.OrderID,
			OrderNo: order.OrderNo, WarehouseID: order.WarehouseID,
			Status: string(order.Status), WaveID: order.WaveID,
			CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
		}
		if err := tx.Create(om).Error; err != nil {
			return err
		}
		for _, item := range order.Items {
			if err := tx.Create(&OutboundItemModel{
				ID: item.ID, OutboundID: order.ID, SKUID: item.SKUID,
				SKUCode: item.SKUCode, SKUName: item.SKUName,
				Quantity: item.Quantity, PickedQty: item.PickedQty,
				CheckedQty: item.CheckedQty, LocationID: item.LocationID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *WarehouseRepository) ListOutbounds(ctx context.Context, tenantID string, offset, limit int) ([]*domain.OutboundOrder, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&OutboundOrderModel{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	var models []*OutboundOrderModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	orders := make([]*domain.OutboundOrder, len(models))
	for i, m := range models {
		var items []*OutboundItemModel
		r.db.WithContext(ctx).Where("outbound_id = ?", m.ID).Find(&items)
		domainItems := make([]*domain.OutboundItem, len(items))
		for j, it := range items {
			domainItems[j] = &domain.OutboundItem{ID: it.ID, SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Quantity: it.Quantity, PickedQty: it.PickedQty, CheckedQty: it.CheckedQty, LocationID: it.LocationID}
		}
		orders[i] = &domain.OutboundOrder{ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, OrderNo: m.OrderNo, WarehouseID: m.WarehouseID, Status: domain.OutboundStatus(m.Status), Items: domainItems, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
	}
	return orders, total, nil
}

func (r *WarehouseRepository) FindOutbound(ctx context.Context, id string) (*domain.OutboundOrder, error) {
	var m OutboundOrderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	var items []*OutboundItemModel
	r.db.WithContext(ctx).Where("outbound_id = ?", m.ID).Find(&items)
	domainItems := make([]*domain.OutboundItem, len(items))
	for i, it := range items {
		domainItems[i] = &domain.OutboundItem{ID: it.ID, SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Quantity: it.Quantity, PickedQty: it.PickedQty, CheckedQty: it.CheckedQty, LocationID: it.LocationID}
	}
	return &domain.OutboundOrder{ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, OrderNo: m.OrderNo, WarehouseID: m.WarehouseID, Status: domain.OutboundStatus(m.Status), Items: domainItems, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (r *WarehouseRepository) UpdateOutboundStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&OutboundOrderModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r *WarehouseRepository) CreatePickTask(ctx context.Context, task *domain.PickTask) error {
	return r.db.WithContext(ctx).Create(&PickTaskModel{
		ID: task.ID, WaveID: task.WaveID, OutboundID: task.OutboundID,
		SKUID: task.SKUID, SKUCode: task.SKUCode, SKUName: task.SKUName,
		Quantity: task.Quantity, PickedQty: task.PickedQty,
		LocationCode: task.LocationCode, Status: task.Status, PickerID: task.PickerID,
	}).Error
}

func (r *WarehouseRepository) ListPickTasks(ctx context.Context, outboundID string) ([]*domain.PickTask, error) {
	var models []*PickTaskModel
	err := r.db.WithContext(ctx).Where("outbound_id = ?", outboundID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	tasks := make([]*domain.PickTask, len(models))
	for i, m := range models {
		tasks[i] = &domain.PickTask{ID: m.ID, WaveID: m.WaveID, OutboundID: m.OutboundID, SKUID: m.SKUID, SKUCode: m.SKUCode, SKUName: m.SKUName, Quantity: m.Quantity, PickedQty: m.PickedQty, LocationCode: m.LocationCode, Status: m.Status, PickerID: m.PickerID}
	}
	return tasks, nil
}

func (r *WarehouseRepository) UpdatePickQty(ctx context.Context, id string, pickedQty int, status string) error {
	return r.db.WithContext(ctx).Model(&PickTaskModel{}).Where("id = ?", id).Updates(map[string]interface{}{"picked_quantity": pickedQty, "status": status}).Error
}

func (r *WarehouseRepository) ListWarehouses(ctx context.Context, tenantID string) ([]*domain.Warehouse, error) {
	var models []*WarehouseModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	whs := make([]*domain.Warehouse, len(models))
	for i, m := range models {
		whs[i] = &domain.Warehouse{ID: m.ID, TenantID: m.TenantID, Name: m.Name, Code: m.Code, Address: m.Address, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return whs, nil
}
