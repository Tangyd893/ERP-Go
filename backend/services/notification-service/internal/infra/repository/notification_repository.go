package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/domain"
	"gorm.io/gorm"
)

const orderByDesc = "created_at DESC"

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Create(&NotificationModel{
		ID: n.ID, TenantID: n.TenantID, UserID: n.UserID,
		Title: n.Title, Content: n.Content, Type: n.Type, Read: n.Read, CreatedAt: n.CreatedAt,
	}).Error
}

func (r *NotificationRepository) ListByUser(ctx context.Context, tenantID, userID string, offset, limit int) ([]*domain.Notification, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID)
	query.Count(&total)
	var models []*NotificationModel
	query.Order(orderByDesc).Offset(offset).Limit(limit).Find(&models)
	list := make([]*domain.Notification, len(models))
	for i, m := range models {
		list[i] = &domain.Notification{ID: m.ID, TenantID: m.TenantID, UserID: m.UserID, Title: m.Title, Content: m.Content, Type: m.Type, Read: m.Read, CreatedAt: m.CreatedAt}
	}
	return list, total, nil
}

func (r *NotificationRepository) CountUnread(ctx context.Context, tenantID, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("tenant_id = ? AND user_id = ? AND read = false", tenantID, userID).Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, tenantID, userID string) error {
	return r.db.WithContext(ctx).Model(&NotificationModel{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Update("read", true).Error
}
