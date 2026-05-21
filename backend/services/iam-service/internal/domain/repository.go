package domain

import (
	"context"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindByID(ctx context.Context, tenantID, userID string) (*User, error)
	FindByUsername(ctx context.Context, tenantID, username string) (*User, error)
	FindWithRoles(ctx context.Context, tenantID, userID string) (*User, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]*User, int64, error)
	Delete(ctx context.Context, tenantID, userID string) error
}

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	FindByID(ctx context.Context, tenantID, roleID string) (*Role, error)
	FindByCode(ctx context.Context, tenantID, code string) (*Role, error)
	FindWithPermissions(ctx context.Context, tenantID, roleID string) (*Role, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]*Role, int64, error)
	Delete(ctx context.Context, tenantID, roleID string) error

	// 角色-权限关联
	AddPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error

	// 用户-角色关联
	AssignUserRoles(ctx context.Context, userID string, roleIDs []string) error
	RemoveUserRoles(ctx context.Context, userID string, roleIDs []string) error
}

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, permission *Permission) error
	FindByID(ctx context.Context, permissionID string) (*Permission, error)
	FindByCode(ctx context.Context, code string) (*Permission, error)
	List(ctx context.Context, offset, limit int) ([]*Permission, int64, error)
	ListByRoleID(ctx context.Context, roleID string) ([]*Permission, error)
	Delete(ctx context.Context, permissionID string) error
}

// AuditRepository 审计仓储接口
type AuditRepository interface {
	Write(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, tenantID string, offset, limit int) ([]*AuditLog, int64, error)
}

// TokenManager Token 管理器接口
type TokenManager interface {
	GenerateAccessToken(userID, tenantID string, roles []Role) (string, error)
	GenerateRefreshToken(userID, tenantID string) (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

// TokenClaims Token 中的声明
type TokenClaims struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// PasswordHasher 密码哈希器接口
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
