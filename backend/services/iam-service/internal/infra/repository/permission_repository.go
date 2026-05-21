package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

// PermissionRepository GORM 实现的权限仓储
type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, perm *domain.Permission) error {
	model := &PermissionModel{
		ID:           perm.ID,
		Name:         perm.Name,
		Code:         perm.Code,
		Description:  perm.Description,
		ResourceType: string(perm.ResourceType),
		Action:       perm.Action,
		ParentID:     perm.ParentID,
		SortOrder:    perm.SortOrder,
		CreatedAt:    perm.CreatedAt,
		UpdatedAt:    perm.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *PermissionRepository) Update(ctx context.Context, perm *domain.Permission) error {
	return r.db.WithContext(ctx).Model(&PermissionModel{}).
		Where("id = ?", perm.ID).
		Updates(map[string]interface{}{
			"name":        perm.Name,
			"description": perm.Description,
			"sort_order":  perm.SortOrder,
		}).Error
}

func (r *PermissionRepository) FindByID(ctx context.Context, permID string) (*domain.Permission, error) {
	var model PermissionModel
	err := r.db.WithContext(ctx).Where("id = ?", permID).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "权限不存在")
		}
		return nil, err
	}
	return modelToDomainPermission(&model), nil
}

func (r *PermissionRepository) FindByCode(ctx context.Context, code string) (*domain.Permission, error) {
	var model PermissionModel
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "权限不存在")
		}
		return nil, err
	}
	return modelToDomainPermission(&model), nil
}

func (r *PermissionRepository) List(ctx context.Context, offset, limit int) ([]*domain.Permission, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&PermissionModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*PermissionModel
	err := r.db.WithContext(ctx).
		Order("sort_order ASC").
		Offset(offset).Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	perms := make([]*domain.Permission, len(models))
	for i, m := range models {
		perms[i] = modelToDomainPermission(m)
	}
	return perms, total, nil
}

func (r *PermissionRepository) ListByRoleID(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	var models []PermissionModel
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ?", roleID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	perms := make([]*domain.Permission, len(models))
	for i, m := range models {
		perms[i] = modelToDomainPermission(&m)
	}
	return perms, nil
}

func (r *PermissionRepository) Delete(ctx context.Context, permID string) error {
	return r.db.WithContext(ctx).Where("id = ?", permID).Delete(&PermissionModel{}).Error
}
