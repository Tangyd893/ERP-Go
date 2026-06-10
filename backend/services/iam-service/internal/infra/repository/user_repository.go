package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

const whereIDAndTenantID = "id = ? AND tenant_id = ?"

// UserRepository GORM 实现的用户仓储
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := &UserModel{
		ID:           user.ID,
		TenantID:     user.TenantID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Status:       string(user.Status),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return errors.WrapError(errors.CodeInternalError, err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	updates := map[string]interface{}{
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"status":     string(user.Status),
		"updated_at": user.UpdatedAt,
	}
	if user.LastLoginAt != nil {
		updates["last_login_at"] = user.LastLoginAt
	}
	return r.db.WithContext(ctx).Model(&UserModel{}).
		Where(whereIDAndTenantID, user.ID, user.TenantID).
		Updates(updates).Error
}

func (r *UserRepository) FindByID(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).
		Where(whereIDAndTenantID, userID, tenantID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
		}
		return nil, err
	}
	return modelToDomainUser(&model), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, tenantID, username string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND username = ?", tenantID, username).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
		}
		return nil, err
	}
	return modelToDomainUser(&model), nil
}

func (r *UserRepository) FindWithRoles(ctx context.Context, tenantID, usernameOrID string) (*domain.User, error) {
	// 先尝试按用户名查找
	user, err := r.FindByUsername(ctx, tenantID, usernameOrID)
	if err != nil {
		// 再尝试按ID查找
		user, err = r.FindByID(ctx, tenantID, usernameOrID)
		if err != nil {
			return nil, err
		}
	}

	// 加载用户角色及其权限
	var roleModels []RoleModel
	err = r.db.WithContext(ctx).
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ? AND roles.tenant_id = ? AND roles.status = ?", user.ID, tenantID, "active").
		Find(&roleModels).Error
	if err != nil {
		return nil, err
	}

	for _, rm := range roleModels {
		role := modelToDomainRole(&rm)

		// 加载角色的权限
		var permModels []PermissionModel
		err = r.db.WithContext(ctx).
			Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
			Where("rp.role_id = ?", rm.ID).
			Find(&permModels).Error
		if err != nil {
			return nil, err
		}
		for _, pm := range permModels {
			role.Permissions = append(role.Permissions, *modelToDomainPermission(&pm))
		}

		user.Roles = append(user.Roles, *role)
	}

	return user, nil
}

func (r *UserRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*UserModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = modelToDomainUser(m)
	}
	return users, total, nil
}

func (r *UserRepository) Delete(ctx context.Context, tenantID, userID string) error {
	return r.db.WithContext(ctx).Where(whereIDAndTenantID, userID, tenantID).Delete(&UserModel{}).Error
}

func modelToDomainUser(m *UserModel) *domain.User {
	return &domain.User{
		ID:           m.ID,
		TenantID:     m.TenantID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Nickname:     m.Nickname,
		Email:        m.Email,
		Phone:        m.Phone,
		Avatar:       m.Avatar,
		Status:       domain.UserStatus(m.Status),
		LastLoginAt:  m.LastLoginAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
