package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/infra/repository"
)

// ChannelAppService 渠道应用服务
type ChannelAppService struct {
	repo *repository.ChannelRepository
}

func NewChannelAppService(repo *repository.ChannelRepository) *ChannelAppService {
	return &ChannelAppService{repo: repo}
}

func (s *ChannelAppService) CreateStore(ctx context.Context, store *domain.Store) error {
	return s.repo.CreateStore(ctx, store)
}

func (s *ChannelAppService) ListStores(ctx context.Context, tenantID string) ([]*domain.Store, error) {
	return s.repo.ListStores(ctx, tenantID)
}

func (s *ChannelAppService) CreateImportTask(ctx context.Context, task *domain.OrderImportTask) error {
	return s.repo.CreateImportTask(ctx, task)
}

func (s *ChannelAppService) GetImportTask(ctx context.Context, idempotencyKey string) (*domain.OrderImportTask, error) {
	return s.repo.FindImportTaskByKey(ctx, idempotencyKey)
}
