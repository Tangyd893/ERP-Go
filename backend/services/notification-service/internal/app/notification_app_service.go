package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/infra/repository"
)

type NotificationAppService struct {
	repo *repository.NotificationRepository
}

func NewNotificationAppService(repo *repository.NotificationRepository) *NotificationAppService {
	return &NotificationAppService{repo: repo}
}

func (s *NotificationAppService) ListNotifications(ctx context.Context, tenantID, userID string, offset, limit int) ([]*domain.Notification, int64, error) {
	return s.repo.ListByUser(ctx, tenantID, userID, offset, limit)
}

func (s *NotificationAppService) GetUnreadCount(ctx context.Context, tenantID, userID string) (int64, error) {
	return s.repo.CountUnread(ctx, tenantID, userID)
}

func (s *NotificationAppService) MarkAllRead(ctx context.Context, tenantID, userID string) error {
	return s.repo.MarkAllRead(ctx, tenantID, userID)
}
