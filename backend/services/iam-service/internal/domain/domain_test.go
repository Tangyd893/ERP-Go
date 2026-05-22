package domain

import (
	"testing"
	"time"
)

func TestUserIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status UserStatus
		want   bool
	}{
		{"激活状态", UserStatusActive, true},
		{"禁用状态", UserStatusDisabled, false},
		{"锁定状态", UserStatusLocked, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Status: tt.status}
			if got := u.IsActive(); got != tt.want {
				t.Errorf("User.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHasPermission(t *testing.T) {
	u := &User{
		Roles: []Role{
			{
				Code: "admin",
				Permissions: []Permission{
					{Code: "order:create"},
					{Code: "order:view"},
				},
			},
			{
				Code: "operator",
				Permissions: []Permission{
					{Code: "order:export"},
				},
			},
		},
	}

	if !u.HasPermission("order:create") {
		t.Error("应具有 order:create 权限")
	}
	if !u.HasPermission("order:view") {
		t.Error("应具有 order:view 权限")
	}
	if !u.HasPermission("order:export") {
		t.Error("应具有 order:export 权限（跨角色）")
	}
	if u.HasPermission("user:delete") {
		t.Error("不应具有 user:delete 权限")
	}

	emptyUser := &User{}
	if emptyUser.HasPermission("any") {
		t.Error("无角色用户不应具有任何权限")
	}
}

func TestUserHasAnyRole(t *testing.T) {
	u := &User{
		Roles: []Role{
			{Code: "admin"},
			{Code: "operator"},
		},
	}

	if !u.HasAnyRole("admin") {
		t.Error("应匹配 admin 角色")
	}
	if !u.HasAnyRole("operator") {
		t.Error("应匹配 operator 角色")
	}
	if !u.HasAnyRole("admin", "operator", "viewer") {
		t.Error("变参中任一角色匹配应返回 true")
	}
	if u.HasAnyRole("viewer") {
		t.Error("不应匹配 viewer 角色")
	}
	if u.HasAnyRole() {
		t.Error("空变参不应匹配任何角色")
	}

	emptyUser := &User{}
	if emptyUser.HasAnyRole("admin") {
		t.Error("无角色用户不应匹配任何角色")
	}
}

func TestUserDisableAndEnable(t *testing.T) {
	u := &User{Status: UserStatusActive}

	u.Disable()
	if u.Status != UserStatusDisabled {
		t.Errorf("Disable 后状态应为 disabled，实际: %s", u.Status)
	}
	if u.IsActive() {
		t.Error("禁用后 IsActive 应为 false")
	}

	u.Enable()
	if u.Status != UserStatusActive {
		t.Errorf("Enable 后状态应为 active，实际: %s", u.Status)
	}
	if !u.IsActive() {
		t.Error("启用后 IsActive 应为 true")
	}
}

func TestUserRecordLogin(t *testing.T) {
	u := &User{}
	if u.LastLoginAt != nil {
		t.Error("新用户的 LastLoginAt 应为 nil")
	}

	before := time.Now()
	u.RecordLogin()
	after := time.Now()

	if u.LastLoginAt == nil {
		t.Fatal("RecordLogin 后 LastLoginAt 不应为 nil")
	}
	if u.LastLoginAt.Before(before) || u.LastLoginAt.After(after) {
		t.Errorf("RecordLogin 时间应处于合理范围: %v (范围: %v ~ %v)", u.LastLoginAt, before, after)
	}
}

func TestRoleIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status RoleStatus
		want   bool
	}{
		{"激活状态", RoleStatusActive, true},
		{"禁用状态", RoleStatusDisabled, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{Status: tt.status}
			if got := r.IsActive(); got != tt.want {
				t.Errorf("Role.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleHasPermission(t *testing.T) {
	r := &Role{
		Permissions: []Permission{
			{Code: "order:create"},
			{Code: "order:view"},
		},
	}

	if !r.HasPermission("order:create") {
		t.Error("应具有 order:create 权限")
	}
	if !r.HasPermission("order:view") {
		t.Error("应具有 order:view 权限")
	}
	if r.HasPermission("order:delete") {
		t.Error("不应具有 order:delete 权限")
	}
	if r.HasPermission("") {
		t.Error("空 Code 不匹配任何权限")
	}

	emptyRole := &Role{}
	if emptyRole.HasPermission("any") {
		t.Error("无权限的角色不应匹配任何权限")
	}
}

func TestRoleAssignPermissions(t *testing.T) {
	r := &Role{Code: "admin"}

	r.AssignPermissions(nil)
	if r.Permissions != nil {
		t.Error("分配 nil 权限后 Permissions 应为 nil")
	}

	perms := []Permission{
		{Code: "order:create"},
		{Code: "order:view"},
	}
	r.AssignPermissions(perms)
	if len(r.Permissions) != 2 {
		t.Errorf("应分配 2 个权限，实际: %d", len(r.Permissions))
	}
	if !r.HasPermission("order:create") {
		t.Error("分配后应具有 order:create 权限")
	}

	r.AssignPermissions([]Permission{})
	if len(r.Permissions) != 0 {
		t.Errorf("重新分配空切片后应为空，实际: %d", len(r.Permissions))
	}
}

func TestPermissionIsMenu(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceType
		want         bool
	}{
		{"菜单类型", ResourceMenu, true},
		{"按钮类型", ResourceButton, false},
		{"API类型", ResourceAPI, false},
		{"数据类型", ResourceData, false},
		{"任意类型", ResourceAny, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Permission{ResourceType: tt.resourceType}
			if got := p.IsMenu(); got != tt.want {
				t.Errorf("Permission.IsMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionIsButton(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceType
		want         bool
	}{
		{"按钮类型", ResourceButton, true},
		{"菜单类型", ResourceMenu, false},
		{"API类型", ResourceAPI, false},
		{"数据类型", ResourceData, false},
		{"任意类型", ResourceAny, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Permission{ResourceType: tt.resourceType}
			if got := p.IsButton(); got != tt.want {
				t.Errorf("Permission.IsButton() = %v, want %v", got, tt.want)
			}
		})
	}
}
