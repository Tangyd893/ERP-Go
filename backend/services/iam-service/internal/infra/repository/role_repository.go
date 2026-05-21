package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

// RoleRepository GORM 实现的角色仓储
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	model := &RoleModel{
		ID:          role.ID,
		TenantID:    role.TenantID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      string(role.Status),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Model(&RoleModel{}).
		Where("id = ? AND tenant_id = ?", role.ID, role.TenantID).
		Updates(map[string]interface{}{
			"name":        role.Name,
			"description": role.Description,
			"status":      string(role.Status),
			"updated_at":  role.UpdatedAt,
		}).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, tenantID, roleID string) (*domain.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", roleID, tenantID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "角色不存在")
		}
		return nil, err
	}
	return modelToDomainRole(&model), nil
}

func (r *RoleRepository) FindByCode(ctx context.Context, tenantID, code string) (*domain.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "角色不存在")
		}
		return nil, err
	}
	return modelToDomainRole(&model), nil
}

func (r *RoleRepository) FindWithPermissions(ctx context.Context, tenantID, roleID string) (*domain.Role, error) {
	role, err := r.FindByID(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}

	var permModels []PermissionModel
	err = r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ?", roleID).
		Find(&permModels).Error
	if err != nil {
		return nil, err
	}

	for _, pm := range permModels {
		role.Permissions = append(role.Permissions, *modelToDomainPermission(&pm))
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Role, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*RoleModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	roles := make([]*domain.Role, len(models))
	for i, m := range models {
		roles[i] = modelToDomainRole(m)
	}
	return roles, total, nil
}

func (r *RoleRepository) Delete(ctx context.Context, tenantID, roleID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色关联的权限
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		// 删除角色关联的用户
		if err := tx.Where("role_id = ?", roleID).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		// 删除角色
		return tx.Where("id = ? AND tenant_id = ?", roleID, tenantID).Delete(&RoleModel{}).Error
	})
}

func (r *RoleRepository) AddPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	records := make([]*RolePermissionModel, len(permissionIDs))
	for i, permID := range permissionIDs {
		records[i] = &RolePermissionModel{
			RoleID:       roleID,
			PermissionID: permID,
		}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *RoleRepository) RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&RolePermissionModel{}).Error
}

func (r *RoleRepository) AssignUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	records := make([]*UserRoleModel, len(roleIDs))
	for i, roleID := range roleIDs {
		records[i] = &UserRoleModel{
			UserID: userID,
			RoleID: roleID,
		}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *RoleRepository) RemoveUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&UserRoleModel{}).Error
}

func modelToDomainRole(m *RoleModel) *domain.Role {
	return &domain.Role{
		ID:          m.ID,
		TenantID:    m.TenantID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Status:      domain.RoleStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func modelToDomainPermission(m *PermissionModel) *domain.Permission {
	return &domain.Permission{
		ID:           m.ID,
		Name:         m.Name,
		Code:         m.Code,
		Description:  m.Description,
		ResourceType: domain.ResourceType(m.ResourceType),
		Action:       m.Action,
		ParentID:     m.ParentID,
		SortOrder:    m.SortOrder,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
