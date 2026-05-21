package domain

import "time"

// File 文件实体
type File struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	Bucket     string    `json:"bucket"`
	ObjectKey  string    `json:"object_key"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	SourceType string    `json:"source_type"`
	SourceID   string    `json:"source_id"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}
