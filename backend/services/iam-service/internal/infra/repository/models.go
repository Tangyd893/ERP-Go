package repository

import "time"

// UserModel 用户表 GORM 模型
type UserModel struct {
	ID           string     `gorm:"column:id;primaryKey"`
	TenantID     string     `gorm:"column:tenant_id;index:idx_users_tenant_username,unique"`
	Username     string     `gorm:"column:username;index:idx_users_tenant_username,unique"`
	PasswordHash string     `gorm:"column:password_hash"`
	Nickname     string     `gorm:"column:nickname"`
	Email        string     `gorm:"column:email"`
	Phone        string     `gorm:"column:phone"`
	Avatar       string     `gorm:"column:avatar"`
	Status       string     `gorm:"column:status;index"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (UserModel) TableName() string { return "users" }

// RoleModel 角色表 GORM 模型
type RoleModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	TenantID    string    `gorm:"column:tenant_id"`
	Name        string    `gorm:"column:name"`
	Code        string    `gorm:"column:code"`
	Description string    `gorm:"column:description"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (RoleModel) TableName() string { return "roles" }

// PermissionModel 权限表 GORM 模型
type PermissionModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	Code         string    `gorm:"column:code;unique"`
	Description  string    `gorm:"column:description"`
	ResourceType string    `gorm:"column:resource_type"`
	Action       string    `gorm:"column:action"`
	ParentID     string    `gorm:"column:parent_id"`
	SortOrder    int       `gorm:"column:sort_order"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (PermissionModel) TableName() string { return "permissions" }

// UserRoleModel 用户-角色关联表 GORM 模型
type UserRoleModel struct {
	UserID    string    `gorm:"column:user_id;primaryKey"`
	RoleID    string    `gorm:"column:role_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserRoleModel) TableName() string { return "user_roles" }

// RolePermissionModel 角色-权限关联表 GORM 模型
type RolePermissionModel struct {
	RoleID       string    `gorm:"column:role_id;primaryKey"`
	PermissionID string    `gorm:"column:permission_id;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }

// AuditLogModel 审计日志表 GORM 模型
type AuditLogModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	TenantID     string    `gorm:"column:tenant_id;index"`
	UserID       string    `gorm:"column:user_id"`
	Username     string    `gorm:"column:username"`
	Action       string    `gorm:"column:action"`
	ResourceType string    `gorm:"column:resource_type"`
	ResourceID   string    `gorm:"column:resource_id"`
	Detail       string    `gorm:"column:detail"`
	IP           string    `gorm:"column:ip"`
	UserAgent    string    `gorm:"column:user_agent"`
	RequestID    string    `gorm:"column:request_id"`
	TraceID      string    `gorm:"column:trace_id"`
	Result       string    `gorm:"column:result"`
	ResultDetail string    `gorm:"column:result_detail"`
	CreatedAt    time.Time `gorm:"column:created_at;index"`
}

func (AuditLogModel) TableName() string { return "audit_logs" }
