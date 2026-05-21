package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/domain"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) CreateStore(ctx context.Context, store *domain.Store) error {
	return r.db.WithContext(ctx).Create(&StoreModel{
		ID: store.ID, TenantID: store.TenantID, Platform: store.PlatformCode,
		Site: store.Site, Name: store.Name, StoreCode: store.StoreCode,
		AuthToken: store.AuthToken, AuthStatus: store.AuthStatus,
		AuthExpiresAt: store.AuthExpiry, Status: string(store.Status),
		CreatedAt: store.CreatedAt, UpdatedAt: store.UpdatedAt,
	}).Error
}

func (r *ChannelRepository) ListStores(ctx context.Context, tenantID string) ([]*domain.Store, error) {
	var models []*StoreModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	stores := make([]*domain.Store, len(models))
	for i, m := range models {
		stores[i] = &domain.Store{
			ID: m.ID, TenantID: m.TenantID, PlatformCode: m.Platform, Site: m.Site,
			Name: m.Name, StoreCode: m.StoreCode, AuthToken: m.AuthToken,
			AuthStatus: m.AuthStatus, AuthExpiry: m.AuthExpiresAt,
			Status: domain.StoreStatus(m.Status), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		}
	}
	return stores, nil
}

func (r *ChannelRepository) CreateImportTask(ctx context.Context, task *domain.OrderImportTask) error {
	return r.db.WithContext(ctx).Create(&OrderImportTaskModel{
		ID: task.ID, TenantID: task.TenantID, StoreID: task.StoreID,
		ImportType: task.ImportType, FileName: task.FileName,
		IdempotencyKey: task.IdempotencyKey, Status: task.Status,
		TotalRows: task.TotalRows, SuccessRows: task.SuccessRows, FailRows: task.FailedRows,
		ErrorMsg: task.ErrorMsg, CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
	}).Error
}

func (r *ChannelRepository) FindImportTaskByKey(ctx context.Context, idempotencyKey string) (*domain.OrderImportTask, error) {
	var m OrderImportTaskModel
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", idempotencyKey).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &domain.OrderImportTask{
		ID: m.ID, TenantID: m.TenantID, StoreID: m.StoreID,
		ImportType: m.ImportType, FileName: m.FileName,
		IdempotencyKey: m.IdempotencyKey, Status: m.Status,
		TotalRows: m.TotalRows, SuccessRows: m.SuccessRows, FailedRows: m.FailRows,
		ErrorMsg: m.ErrorMsg, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}
