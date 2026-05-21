package repository

import "time"

// TenantModel 租户表 GORM 模型
type TenantModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	Code         string    `gorm:"column:code;unique"`
	ContactName  string    `gorm:"column:contact_name"`
	ContactEmail string    `gorm:"column:contact_email"`
	ContactPhone string    `gorm:"column:contact_phone"`
	Status       string    `gorm:"column:status"`
	QuotaUsers   int       `gorm:"column:quota_users"`
	QuotaOrders  int       `gorm:"column:quota_orders"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (TenantModel) TableName() string { return "tenants" }

// OrganizationModel 组织表 GORM 模型
type OrganizationModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	ParentID  string    `gorm:"column:parent_id"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	SortOrder int       `gorm:"column:sort_order"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (OrganizationModel) TableName() string { return "organizations" }

// DepartmentModel 部门表 GORM 模型
type DepartmentModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	OrgID     string    `gorm:"column:org_id"`
	ParentID  string    `gorm:"column:parent_id"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	ManagerID string    `gorm:"column:manager_id"`
	SortOrder int       `gorm:"column:sort_order"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (DepartmentModel) TableName() string { return "departments" }

// PositionModel 岗位表 GORM 模型
type PositionModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	DeptID    string    `gorm:"column:dept_id"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	SortOrder int       `gorm:"column:sort_order"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PositionModel) TableName() string { return "positions" }
