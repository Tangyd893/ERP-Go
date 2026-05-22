package domain

import "time"

// IsActive 租户是否激活
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

// IsSuspended 租户是否暂停
func (t *Tenant) IsSuspended() bool {
	return t.Status == TenantStatusSuspended
}

// IsDisabled 租户是否禁用
func (t *Tenant) IsDisabled() bool {
	return t.Status == TenantStatusDisabled
}

// Disable 禁用租户
func (t *Tenant) Disable() {
	t.Status = TenantStatusDisabled
	t.UpdatedAt = time.Now()
}

// Suspend 暂停租户
func (t *Tenant) Suspend() {
	t.Status = TenantStatusSuspended
	t.UpdatedAt = time.Now()
}
