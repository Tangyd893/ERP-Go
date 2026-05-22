package domain

import "testing"

func TestTenantIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status TenantStatus
		want   bool
	}{
		{"激活状态", TenantStatusActive, true},
		{"禁用状态", TenantStatusDisabled, false},
		{"暂停状态", TenantStatusSuspended, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tenant{Status: tt.status}
			if got := tr.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantIsSuspended(t *testing.T) {
	tests := []struct {
		name   string
		status TenantStatus
		want   bool
	}{
		{"暂停状态", TenantStatusSuspended, true},
		{"激活状态", TenantStatusActive, false},
		{"禁用状态", TenantStatusDisabled, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tenant{Status: tt.status}
			if got := tr.IsSuspended(); got != tt.want {
				t.Errorf("IsSuspended() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantDisableAndSuspend(t *testing.T) {
	tr := &Tenant{Status: TenantStatusActive}

	tr.Disable()
	if tr.Status != TenantStatusDisabled {
		t.Errorf("禁用后状态应为 disabled，实际: %s", tr.Status)
	}
	if tr.IsActive() {
		t.Error("禁用后 IsActive 应为 false")
	}

	tr2 := &Tenant{Status: TenantStatusActive}
	tr2.Suspend()
	if tr2.Status != TenantStatusSuspended {
		t.Errorf("暂停后状态应为 suspended，实际: %s", tr2.Status)
	}
	if !tr2.IsSuspended() {
		t.Error("暂停后 IsSuspended 应为 true")
	}
}
