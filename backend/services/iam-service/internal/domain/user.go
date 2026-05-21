package domain

import "time"

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusDisabled  UserStatus = "disabled"
	UserStatusLocked    UserStatus = "locked"
)

// User 用户聚合根
type User struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	Nickname     string     `json:"nickname"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Avatar       string     `json:"avatar"`
	Status       UserStatus `json:"status"`
	Roles        []Role     `json:"roles,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// IsActive 用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// HasPermission 检查用户是否有某权限
func (u *User) HasPermission(code string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Code == code {
				return true
			}
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任一角色
func (u *User) HasAnyRole(codes ...string) bool {
	for _, role := range u.Roles {
		for _, code := range codes {
			if role.Code == code {
				return true
			}
		}
	}
	return false
}

// Disable 禁用用户
func (u *User) Disable() {
	u.Status = UserStatusDisabled
}

// Enable 启用用户
func (u *User) Enable() {
	u.Status = UserStatusActive
}

// RecordLogin 记录登录时间
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}
