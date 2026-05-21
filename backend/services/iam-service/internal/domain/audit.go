package domain

import "time"

// AuditAction 审计操作类型
type AuditAction string

const (
	AuditLogin        AuditAction = "login"
	AuditLogout       AuditAction = "logout"
	AuditCreate       AuditAction = "create"
	AuditUpdate       AuditAction = "update"
	AuditDelete       AuditAction = "delete"
	AuditExport       AuditAction = "export"
	AuditImport       AuditAction = "import"
	AuditPermissionChange AuditAction = "permission_change"
	AuditDataAccess   AuditAction = "data_access"
)

// AuditLog 操作审计日志实体
type AuditLog struct {
	ID           string      `json:"id"`
	TenantID     string      `json:"tenant_id"`
	UserID       string      `json:"user_id"`
	Username     string      `json:"username"`
	Action       AuditAction `json:"action"`
	ResourceType string      `json:"resource_type"`
	ResourceID   string      `json:"resource_id,omitempty"`
	Detail       string      `json:"detail,omitempty"`
	IP           string      `json:"ip"`
	UserAgent    string      `json:"user_agent,omitempty"`
	RequestID    string      `json:"request_id,omitempty"`
	TraceID      string      `json:"trace_id,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

// AuditResult 审计结果信息
type AuditResult struct {
	Result  string `json:"result"`
	Details string `json:"details,omitempty"`
}
