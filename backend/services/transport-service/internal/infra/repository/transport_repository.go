package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/domain"
	"gorm.io/gorm"
)

const (
	whereTenantID = "tenant_id = ?"
	whereID       = "id = ?"
	orderByDesc   = "created_at DESC"
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
	query := r.db.WithContext(ctx).Model(&ShipmentModel{}).Where(whereTenantID, tenantID)
	query.Count(&total)
	var models []*ShipmentModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
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
	if err := r.db.WithContext(ctx).Where(whereID, id).First(&m).Error; err != nil {
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
	return r.db.WithContext(ctx).Model(&ShipmentModel{}).Where(whereID, id).Updates(updates).Error
}

// ListShippingRules 按租户查询物流规则（按优先级升序）
func (r *TransportRepository) ListShippingRules(ctx context.Context, tenantID string) ([]*domain.ShippingRule, error) {
	var models []*ShippingRuleModel
	err := r.db.WithContext(ctx).Where(whereTenantID, tenantID).
		Order("priority ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	rules := make([]*domain.ShippingRule, len(models))
	for i, m := range models {
		var codes []string
		if m.CountryCodes != "" {
			codes = parseJSONArray(m.CountryCodes)
		}
		rules[i] = &domain.ShippingRule{
			ID: m.ID, TenantID: m.TenantID, Name: m.Name, Priority: m.Priority,
			CountryCodes: codes, MinWeight: m.MinWeight, MaxWeight: m.MaxWeight,
			CarrierServiceID: m.CarrierServiceID,
		}
	}
	return rules, nil
}

// ListCarrierServices 按物流商查询物流产品
func (r *TransportRepository) ListCarrierServices(ctx context.Context, carrierID string) ([]*domain.CarrierService, error) {
	var models []*CarrierServiceModel
	err := r.db.WithContext(ctx).Where("carrier_id = ?", carrierID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	services := make([]*domain.CarrierService, len(models))
	for i, m := range models {
		services[i] = &domain.CarrierService{
			ID: m.ID, CarrierID: m.CarrierID, Name: m.Name, Code: m.Code, ServiceType: m.ServiceType,
		}
	}
	return services, nil
}

// FindCarrierService 按 ID 查询物流产品
func (r *TransportRepository) FindCarrierService(ctx context.Context, id string) (*domain.CarrierService, error) {
	var m CarrierServiceModel
	if err := r.db.WithContext(ctx).Where(whereID, id).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.CarrierService{
		ID: m.ID, CarrierID: m.CarrierID, Name: m.Name, Code: m.Code, ServiceType: m.ServiceType,
	}, nil
}

// parseJSONArray 解析 JSON 字符串数组
func parseJSONArray(s string) []string {
	var arr []string
	// 简单解析 ["a","b"] 格式
	inQuote := false
	cur := ""
	for _, c := range s {
		if c == '"' {
			inQuote = !inQuote
			if !inQuote && cur != "" {
				arr = append(arr, cur)
				cur = ""
			}
			continue
		}
		if inQuote {
			cur += string(c)
		}
	}
	return arr
}

func (r *TransportRepository) ListCarriers(ctx context.Context, tenantID string) ([]*domain.Carrier, error) {
	var models []*CarrierModel
	err := r.db.WithContext(ctx).Where(whereTenantID, tenantID).Find(&models).Error
	if err != nil { return nil, err }
	carriers := make([]*domain.Carrier, len(models))
	for i, m := range models {
		carriers[i] = &domain.Carrier{ID: m.ID, TenantID: m.TenantID, Name: m.Name, Code: m.Code, Status: m.Status, CreatedAt: m.CreatedAt}
	}
	return carriers, nil
}
