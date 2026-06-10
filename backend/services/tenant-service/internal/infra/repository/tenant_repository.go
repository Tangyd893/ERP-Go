package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

const orderByDesc = "created_at DESC"

// TenantRepository GORM 实现的租户仓储
type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	model := &TenantModel{
		ID:           tenant.ID,
		Name:         tenant.Name,
		Code:         tenant.Code,
		ContactName:  tenant.ContactName,
		ContactEmail: tenant.ContactEmail,
		ContactPhone: tenant.ContactPhone,
		Status:       string(tenant.Status),
		QuotaUsers:   tenant.QuotaUsers,
		QuotaOrders:  tenant.QuotaOrders,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Model(&TenantModel{}).
		Where("id = ?", tenant.ID).
		Updates(map[string]interface{}{
			"name":          tenant.Name,
			"contact_name":  tenant.ContactName,
			"contact_email": tenant.ContactEmail,
			"contact_phone": tenant.ContactPhone,
			"status":        string(tenant.Status),
			"quota_users":   tenant.QuotaUsers,
			"quota_orders":  tenant.QuotaOrders,
			"updated_at":    tenant.UpdatedAt,
		}).Error
}

func (r *TenantRepository) FindByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var model TenantModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeTenantNotFound, "租户不存在")
		}
		return nil, err
	}
	return modelToDomainTenant(&model), nil
}

func (r *TenantRepository) FindByCode(ctx context.Context, code string) (*domain.Tenant, error) {
	var model TenantModel
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeTenantNotFound, "租户不存在")
		}
		return nil, err
	}
	return modelToDomainTenant(&model), nil
}

func (r *TenantRepository) List(ctx context.Context, offset, limit int) ([]*domain.Tenant, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&TenantModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*TenantModel
	err := r.db.WithContext(ctx).Order(orderByDesc).Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	tenants := make([]*domain.Tenant, len(models))
	for i, m := range models {
		tenants[i] = modelToDomainTenant(m)
	}
	return tenants, total, nil
}

func modelToDomainTenant(m *TenantModel) *domain.Tenant {
	return &domain.Tenant{
		ID:           m.ID,
		Name:         m.Name,
		Code:         m.Code,
		ContactName:  m.ContactName,
		ContactEmail: m.ContactEmail,
		ContactPhone: m.ContactPhone,
		Status:       domain.TenantStatus(m.Status),
		QuotaUsers:   m.QuotaUsers,
		QuotaOrders:  m.QuotaOrders,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
