package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/domain"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, f *domain.File) error {
	return r.db.WithContext(ctx).Create(&FileModel{
		ID: f.ID, TenantID: f.TenantID, Bucket: f.Bucket, ObjectKey: f.ObjectKey,
		FileName: f.FileName, FileSize: f.FileSize, MimeType: f.MimeType,
		SourceType: f.SourceType, SourceID: f.SourceID, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt,
	}).Error
}

func (r *FileRepository) FindByID(ctx context.Context, id string) (*domain.File, error) {
	var m FileModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.File{ID: m.ID, TenantID: m.TenantID, Bucket: m.Bucket, ObjectKey: m.ObjectKey, FileName: m.FileName, FileSize: m.FileSize, MimeType: m.MimeType, SourceType: m.SourceType, SourceID: m.SourceID, CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt}, nil
}
