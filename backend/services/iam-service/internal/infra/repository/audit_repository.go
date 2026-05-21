package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"gorm.io/gorm"
)

// AuditRepository GORM 实现的审计仓储
type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Write(ctx context.Context, log *domain.AuditLog) error {
	model := &AuditLogModel{
		ID:           log.ID,
		TenantID:     log.TenantID,
		UserID:       log.UserID,
		Username:     log.Username,
		Action:       string(log.Action),
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Detail:       log.Detail,
		IP:           log.IP,
		UserAgent:    log.UserAgent,
		RequestID:    log.RequestID,
		TraceID:      log.TraceID,
		Result:       log.Result,
		ResultDetail: log.ResultDetail,
		CreatedAt:    log.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *AuditRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&AuditLogModel{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*AuditLogModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	logs := make([]*domain.AuditLog, len(models))
	for i, m := range models {
		logs[i] = &domain.AuditLog{
			ID:           m.ID,
			TenantID:     m.TenantID,
			UserID:       m.UserID,
			Username:     m.Username,
			Action:       domain.AuditAction(m.Action),
			ResourceType: m.ResourceType,
			ResourceID:   m.ResourceID,
			Detail:       m.Detail,
			IP:           m.IP,
			UserAgent:    m.UserAgent,
			RequestID:    m.RequestID,
			TraceID:      m.TraceID,
			CreatedAt:    m.CreatedAt,
		}
	}
	return logs, total, nil
}
