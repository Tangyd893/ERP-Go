package domain

import "time"

// RoleStatus 角色状态
type RoleStatus string

const (
	RoleStatusActive   RoleStatus = "active"
	RoleStatusDisabled RoleStatus = "disabled"
)

// Role 角色实体
type Role struct {
	ID          string       `json:"id"`
	TenantID    string       `json:"tenant_id"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Status      RoleStatus   `json:"status"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// IsActive 角色是否激活
func (r *Role) IsActive() bool {
	return r.Status == RoleStatusActive
}

// HasPermission 角色是否有某权限
func (r *Role) HasPermission(code string) bool {
	for _, p := range r.Permissions {
		if p.Code == code {
			return true
		}
	}
	return false
}

// AssignPermissions 分配权限
func (r *Role) AssignPermissions(permissions []Permission) {
	r.Permissions = permissions
}
