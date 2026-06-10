package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/domain"
	"gorm.io/gorm"
)

const orderSortASC = "sort_order ASC"

// OrgRepository GORM 实现的组织/部门/岗位仓储
type OrgRepository struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) *OrgRepository {
	return &OrgRepository{db: db}
}

// 组织相关

func (r *OrgRepository) CreateOrg(ctx context.Context, org *domain.Organization) error {
	model := &OrganizationModel{
		ID:        org.ID,
		TenantID:  org.TenantID,
		ParentID:  org.ParentID,
		Name:      org.Name,
		Code:      org.Code,
		SortOrder: org.SortOrder,
		Status:    org.Status,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *OrgRepository) ListOrgsByTenant(ctx context.Context, tenantID string) ([]*domain.Organization, error) {
	var models []*OrganizationModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order(orderSortASC).Find(&models).Error
	if err != nil {
		return nil, err
	}
	orgs := make([]*domain.Organization, len(models))
	for i, m := range models {
		orgs[i] = modelToDomainOrg(m)
	}
	return orgs, nil
}

// 部门相关

func (r *OrgRepository) CreateDept(ctx context.Context, dept *domain.Department) error {
	model := &DepartmentModel{
		ID:        dept.ID,
		TenantID:  dept.TenantID,
		OrgID:     dept.OrgID,
		ParentID:  dept.ParentID,
		Name:      dept.Name,
		Code:      dept.Code,
		ManagerID: dept.ManagerID,
		SortOrder: dept.SortOrder,
		Status:    dept.Status,
		CreatedAt: dept.CreatedAt,
		UpdatedAt: dept.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *OrgRepository) ListDeptsByOrg(ctx context.Context, orgID string) ([]*domain.Department, error) {
	var models []*DepartmentModel
	err := r.db.WithContext(ctx).Where("org_id = ?", orgID).Order(orderSortASC).Find(&models).Error
	if err != nil {
		return nil, err
	}
	depts := make([]*domain.Department, len(models))
	for i, m := range models {
		depts[i] = modelToDomainDept(m)
	}
	return depts, nil
}

// 岗位相关

func (r *OrgRepository) CreatePosition(ctx context.Context, pos *domain.Position) error {
	model := &PositionModel{
		ID:        pos.ID,
		TenantID:  pos.TenantID,
		DeptID:    pos.DeptID,
		Name:      pos.Name,
		Code:      pos.Code,
		SortOrder: pos.SortOrder,
		Status:    pos.Status,
		CreatedAt: pos.CreatedAt,
		UpdatedAt: pos.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *OrgRepository) ListPositionsByDept(ctx context.Context, deptID string) ([]*domain.Position, error) {
	var models []*PositionModel
	err := r.db.WithContext(ctx).Where("dept_id = ?", deptID).Order(orderSortASC).Find(&models).Error
	if err != nil {
		return nil, err
	}
	positions := make([]*domain.Position, len(models))
	for i, m := range models {
		positions[i] = modelToDomainPosition(m)
	}
	return positions, nil
}

func modelToDomainOrg(m *OrganizationModel) *domain.Organization {
	return &domain.Organization{
		ID:        m.ID,
		TenantID:  m.TenantID,
		ParentID:  m.ParentID,
		Name:      m.Name,
		Code:      m.Code,
		SortOrder: m.SortOrder,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func modelToDomainDept(m *DepartmentModel) *domain.Department {
	return &domain.Department{
		ID:        m.ID,
		TenantID:  m.TenantID,
		OrgID:     m.OrgID,
		ParentID:  m.ParentID,
		Name:      m.Name,
		Code:      m.Code,
		ManagerID: m.ManagerID,
		SortOrder: m.SortOrder,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func modelToDomainPosition(m *PositionModel) *domain.Position {
	return &domain.Position{
		ID:        m.ID,
		TenantID:  m.TenantID,
		DeptID:    m.DeptID,
		Name:      m.Name,
		Code:      m.Code,
		SortOrder: m.SortOrder,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
