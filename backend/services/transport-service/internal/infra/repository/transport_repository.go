package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/domain"
	"gorm.io/gorm"
)

type TransportRepository struct {
	db *gorm.DB
}

func NewTransportRepository(db *gorm.DB) *TransportRepository {
	return &TransportRepository{db: db}
}

func (r *TransportRepository) CreateShipment(ctx context.Context, s *domain.Shipment) error {
	return r.db.WithContext(ctx).Create(&ShipmentModel{
		ID: s.ID, TenantID: s.TenantID, OrderID: s.OrderID, OutboundID: s.OutboundID,
		CarrierCode: s.CarrierCode, ServiceCode: s.ServiceCode, TrackingNo: s.TrackingNo,
		LabelURL: s.LabelURL, Status: string(s.Status), Weight: s.Weight,
		ShippingCost: s.ShippingCost, Currency: s.Currency, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}).Error
}

func (r *TransportRepository) ListShipments(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Shipment, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&ShipmentModel{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	var models []*ShipmentModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	shipments := make([]*domain.Shipment, len(models))
	for i, m := range models {
		shipments[i] = &domain.Shipment{
			ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, OutboundID: m.OutboundID,
			CarrierCode: m.CarrierCode, ServiceCode: m.ServiceCode, TrackingNo: m.TrackingNo,
			LabelURL: m.LabelURL, Status: domain.ShipmentStatus(m.Status), Weight: m.Weight,
			ShippingCost: m.ShippingCost, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		}
	}
	return shipments, total, nil
}

func (r *TransportRepository) FindShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	var m ShipmentModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Shipment{
		ID: m.ID, TenantID: m.TenantID, OrderID: m.OrderID, OutboundID: m.OutboundID,
		CarrierCode: m.CarrierCode, ServiceCode: m.ServiceCode, TrackingNo: m.TrackingNo,
		LabelURL: m.LabelURL, Status: domain.ShipmentStatus(m.Status), Weight: m.Weight,
		ShippingCost: m.ShippingCost, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *TransportRepository) UpdateShipmentStatus(ctx context.Context, id, status string, trackingNo string) error {
	updates := map[string]interface{}{"status": status}
	if trackingNo != "" { updates["tracking_no"] = trackingNo }
	return r.db.WithContext(ctx).Model(&ShipmentModel{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TransportRepository) ListCarriers(ctx context.Context, tenantID string) ([]*domain.Carrier, error) {
	var models []*CarrierModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&models).Error
	if err != nil { return nil, err }
	carriers := make([]*domain.Carrier, len(models))
	for i, m := range models {
		carriers[i] = &domain.Carrier{ID: m.ID, TenantID: m.TenantID, Name: m.Name, Code: m.Code, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return carriers, nil
}
