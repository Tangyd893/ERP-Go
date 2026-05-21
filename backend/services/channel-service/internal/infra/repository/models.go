package repository

import "time"

type StoreModel struct {
	ID            string     `gorm:"column:id;primaryKey"`
	TenantID      string     `gorm:"column:tenant_id;index"`
	Platform      string     `gorm:"column:platform"`
	Site          string     `gorm:"column:site"`
	Name          string     `gorm:"column:name"`
	StoreCode     string     `gorm:"column:store_code"`
	AuthToken     string     `gorm:"column:auth_token"`
	AuthStatus    string     `gorm:"column:auth_status"`
	AuthExpiresAt *time.Time `gorm:"column:auth_expires_at"`
	Status        string     `gorm:"column:status"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (StoreModel) TableName() string { return "stores" }

type SyncTaskModel struct {
	ID           string     `gorm:"column:id;primaryKey"`
	TenantID     string     `gorm:"column:tenant_id"`
	StoreID      string     `gorm:"column:store_id;index"`
	TaskType     string     `gorm:"column:task_type"`
	Status       string     `gorm:"column:status"`
	TotalCount   int        `gorm:"column:total_count"`
	SuccessCount int        `gorm:"column:success_count"`
	FailCount    int        `gorm:"column:fail_count"`
	ErrorMsg     string     `gorm:"column:error_msg"`
	StartedAt    *time.Time `gorm:"column:started_at"`
	EndedAt      *time.Time `gorm:"column:ended_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
}

func (SyncTaskModel) TableName() string { return "sync_tasks" }

type OrderImportTaskModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TenantID       string    `gorm:"column:tenant_id"`
	StoreID        string    `gorm:"column:store_id;index"`
	ImportType     string    `gorm:"column:import_type"`
	FileName       string    `gorm:"column:file_name"`
	IdempotencyKey string    `gorm:"column:idempotency_key;unique"`
	Status         string    `gorm:"column:status"`
	TotalRows      int       `gorm:"column:total_rows"`
	SuccessRows    int       `gorm:"column:success_rows"`
	FailRows       int       `gorm:"column:fail_rows"`
	ErrorMsg       string    `gorm:"column:error_msg"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (OrderImportTaskModel) TableName() string { return "order_import_tasks" }
