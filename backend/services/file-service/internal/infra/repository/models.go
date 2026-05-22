package repository

import "time"

type FileModel struct {
	ID         string    `gorm:"column:id;primaryKey"`
	TenantID   string    `gorm:"column:tenant_id;index"`
	Bucket     string    `gorm:"column:bucket"`
	ObjectKey  string    `gorm:"column:object_key"`
	FileName   string    `gorm:"column:file_name"`
	FileSize   int64     `gorm:"column:file_size"`
	MimeType   string    `gorm:"column:mime_type"`
	SourceType string    `gorm:"column:source_type"`
	SourceID   string    `gorm:"column:source_id"`
	CreatedBy  string    `gorm:"column:created_by"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}
func (FileModel) TableName() string { return "files" }
