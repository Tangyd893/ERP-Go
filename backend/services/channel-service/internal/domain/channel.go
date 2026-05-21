package domain

import "time"

// StoreStatus 店铺状态
type StoreStatus string

const (
	StoreStatusActive     StoreStatus = "active"
	StoreStatusExpired    StoreStatus = "expired"
	StoreStatusSuspended  StoreStatus = "suspended"
)

// Store 店铺聚合根
type Store struct {
	ID          string      `json:"id"`
	TenantID    string      `json:"tenant_id"`
	PlatformCode string     `json:"platform_code"`
	Site        string      `json:"site"`
	Name        string      `json:"name"`
	StoreCode   string      `json:"store_code"`
	AuthToken   string      `json:"-"` // 加密存储
	AuthStatus  string      `json:"auth_status"`
	AuthExpiry  *time.Time  `json:"auth_expiry,omitempty"`
	Status      StoreStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// SyncTask 同步任务
type SyncTask struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	StoreID    string    `json:"store_id"`
	TaskType   string    `json:"task_type"` // order_sync, inventory_push, tracking_upload
	Status     string    `json:"status"`    // pending, running, completed, failed
	TotalCount int       `json:"total_count"`
	SuccessCnt int       `json:"success_count"`
	FailedCnt  int       `json:"failed_count"`
	ErrorMsg   string    `json:"error_msg,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// PlatformAPILog 平台 API 调用日志
type PlatformAPILog struct {
	ID          string    `json:"id"`
	StoreID     string    `json:"store_id"`
	Action      string    `json:"action"`
	RequestURL  string    `json:"request_url"`
	RequestBody string    `json:"request_body,omitempty"`
	StatusCode  int       `json:"status_code"`
	ResponseBody string   `json:"response_body,omitempty"`
	Duration    int64     `json:"duration_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderImportTask 订单导入任务（聚合根）
type OrderImportTask struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	StoreID     string    `json:"store_id"`
	ImportType  string    `json:"import_type"` // csv, api, manual
	FileName    string    `json:"file_name,omitempty"`
	IdempotencyKey string `json:"idempotency_key"`
	Status      string    `json:"status"`    // pending, processing, completed, failed
	TotalRows   int       `json:"total_rows"`
	SuccessRows int       `json:"success_rows"`
	FailedRows  int       `json:"failed_rows"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
