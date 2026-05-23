package domain

import (
	"testing"
	"time"
)

// 创建测试用文件实体
func setupFile() *File {
	return &File{
		ID:         "file-001",
		TenantID:   "tenant-001",
		Bucket:     "erp-uploads",
		ObjectKey:  "orders/2024/01/import_20240115.csv",
		FileName:   "import_20240115.csv",
		FileSize:   204800, // 200KB
		MimeType:   "text/csv",
		SourceType: "order_import",
		SourceID:   "import-001",
		CreatedBy:  "user-001",
		CreatedAt:  time.Now(),
	}
}

// TestFileCreation 测试文件实体创建与基础字段
func TestFileCreation(t *testing.T) {
	f := setupFile()

	if f.ID == "" {
		t.Error("文件ID不应为空")
	}
	if f.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if f.Bucket == "" {
		t.Error("存储桶不应为空")
	}
	if f.ObjectKey == "" {
		t.Error("对象键不应为空")
	}
	if f.FileName == "" {
		t.Error("文件名不应为空")
	}
	if f.FileSize <= 0 {
		t.Error("文件大小应大于0")
	}
	if f.MimeType == "" {
		t.Error("MIME类型不应为空")
	}
	if f.CreatedBy == "" {
		t.Error("创建者不应为空")
	}
}

// TestFileSizeBoundary 测试文件大小边界
func TestFileSizeBoundary(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
	}{
		{"空文件", 0},
		{"1字节", 1},
		{"1KB", 1024},
		{"1MB", 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
		{"100MB", 100 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := setupFile()
			f.FileSize = tt.fileSize
			if f.FileSize != tt.fileSize {
				t.Errorf("文件大小应为 %d，实际 %d", tt.fileSize, f.FileSize)
			}
		})
	}
}

// TestFileMimeType 测试不同 MIME 类型
func TestFileMimeType(t *testing.T) {
	mimeTypes := []string{
		"text/csv",
		"application/pdf",
		"image/png",
		"image/jpeg",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/zip",
		"application/json",
	}

	for _, mime := range mimeTypes {
		t.Run("MIME_"+mime, func(t *testing.T) {
			f := setupFile()
			f.MimeType = mime
			if f.MimeType != mime {
				t.Errorf("MIME类型应为 %s，实际 %s", mime, f.MimeType)
			}
		})
	}
}

// TestFileSourceType 测试不同来源类型
func TestFileSourceType(t *testing.T) {
	sourceTypes := []string{
		"order_import",
		"product_image",
		"shipping_label",
		"invoice",
		"report_export",
		"contract",
	}

	for _, st := range sourceTypes {
		t.Run("来源_"+st, func(t *testing.T) {
			f := setupFile()
			f.SourceType = st
			if f.SourceType != st {
				t.Errorf("来源类型应为 %s，实际 %s", st, f.SourceType)
			}
		})
	}
}

// TestFileSourceID 测试关联来源ID
func TestFileSourceID(t *testing.T) {
	f := setupFile()

	if f.SourceID == "" {
		t.Error("来源ID不应为空")
	}

	// 不同来源应有不同来源ID
	f2 := setupFile()
	f2.SourceType = "product_image"
	f2.SourceID = "prod-001"

	if f.SourceID == f2.SourceID {
		t.Error("不同来源应有不同SourceID")
	}
}

// TestFileBucket 测试不同存储桶
func TestFileBucket(t *testing.T) {
	buckets := []string{"erp-uploads", "erp-exports", "erp-temp", "erp-archive"}

	for _, bucket := range buckets {
		t.Run("桶_"+bucket, func(t *testing.T) {
			f := setupFile()
			f.Bucket = bucket
			if f.Bucket != bucket {
				t.Errorf("存储桶应为 %s，实际 %s", bucket, f.Bucket)
			}
		})
	}
}

// TestFileObjectKey 测试对象键路径
func TestFileObjectKey(t *testing.T) {
	tests := []struct {
		name      string
		objectKey string
	}{
		{"订单导入路径", "orders/2024/01/import_20240115.csv"},
		{"商品图片路径", "products/images/sku-001/main.jpg"},
		{"运单标签路径", "shipping/labels/label-001.pdf"},
		{"报表导出路径", "reports/2024-Q1/sales_summary.xlsx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := setupFile()
			f.ObjectKey = tt.objectKey
			if f.ObjectKey != tt.objectKey {
				t.Errorf("对象键应为 %s，实际 %s", tt.objectKey, f.ObjectKey)
			}
		})
	}
}

// TestFileWithoutSource 测试无关联来源的文件
func TestFileWithoutSource(t *testing.T) {
	f := &File{
		ID:         "file-no-source",
		TenantID:   "tenant-001",
		Bucket:     "erp-uploads",
		ObjectKey:  "misc/standalone_doc.pdf",
		FileName:   "standalone_doc.pdf",
		FileSize:   50000,
		MimeType:   "application/pdf",
		SourceType: "",
		SourceID:   "",
		CreatedBy:  "user-001",
		CreatedAt:  time.Now(),
	}

	if f.SourceType != "" || f.SourceID != "" {
		t.Error("无来源文件 SourceType 和 SourceID 应为空")
	}
}

// TestFileCreatedAt 测试创建时间
func TestFileCreatedAt(t *testing.T) {
	before := time.Now().Add(-1 * time.Minute)
	f := &File{
		ID:         "file-time",
		TenantID:   "tenant-001",
		Bucket:     "erp-uploads",
		ObjectKey:  "test.txt",
		FileName:   "test.txt",
		FileSize:   100,
		MimeType:   "text/plain",
		SourceType: "test",
		SourceID:   "test-001",
		CreatedBy:  "user-001",
		CreatedAt:  time.Now(),
	}
	after := time.Now().Add(1 * time.Minute)

	if f.CreatedAt.Before(before) {
		t.Error("创建时间不应早于测试开始")
	}
	if f.CreatedAt.After(after) {
		t.Error("创建时间不应晚于测试结束")
	}
}

// TestFileLargeSize 测试大文件
func TestFileLargeSize(t *testing.T) {
	f := setupFile()
	f.FileSize = 500 * 1024 * 1024 // 500MB

	if f.FileSize != 500*1024*1024 {
		t.Errorf("大文件大小应为 %d，实际 %d", 500*1024*1024, f.FileSize)
	}
}

// TestFileCreatedBy 测试不同创建者
func TestFileCreatedBy(t *testing.T) {
	users := []string{"user-001", "user-002", "admin", "system"}

	for _, u := range users {
		t.Run("创建者_"+u, func(t *testing.T) {
			f := setupFile()
			f.CreatedBy = u
			if f.CreatedBy != u {
				t.Errorf("创建者应为 %s，实际 %s", u, f.CreatedBy)
			}
		})
	}
}
