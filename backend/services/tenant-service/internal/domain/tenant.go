package domain

import "time"

// TenantStatus 租户状态
type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusDisabled TenantStatus = "disabled"
	TenantStatusSuspended TenantStatus = "suspended"
)

// Tenant 租户聚合根
type Tenant struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	ContactName string       `json:"contact_name"`
	ContactEmail string      `json:"contact_email"`
	ContactPhone string      `json:"contact_phone"`
	Status      TenantStatus `json:"status"`
	QuotaUsers  int          `json:"quota_users"`
	QuotaOrders int          `json:"quota_orders"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Organization 组织实体
type Organization struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	ParentID  string    `json:"parent_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	SortOrder int       `json:"sort_order"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Department 部门实体
type Department struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	OrgID     string    `json:"org_id"`
	ParentID  string    `json:"parent_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	ManagerID string    `json:"manager_id"`
	SortOrder int       `json:"sort_order"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Position 岗位实体
type Position struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	DeptID    string    `json:"dept_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	SortOrder int       `json:"sort_order"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
