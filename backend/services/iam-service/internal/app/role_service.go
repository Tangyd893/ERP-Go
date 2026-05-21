package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/google/uuid"
)

// RoleService 角色管理应用服务
type RoleService struct {
	roleRepo  domain.RoleRepository
	permRepo  domain.PermissionRepository
	auditRepo domain.AuditRepository
}

// NewRoleService 创建角色服务
func NewRoleService(
	roleRepo domain.RoleRepository,
	permRepo domain.PermissionRepository,
	auditRepo domain.AuditRepository,
) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		permRepo:  permRepo,
		auditRepo: auditRepo,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, tenantID, name, code, description string) (*domain.Role, error) {
	existing, _ := s.roleRepo.FindByCode(ctx, tenantID, code)
	if existing != nil {
		return nil, errors.NewBusinessError(errors.CodeAlreadyExists, "角色编码已存在")
	}

	role := &domain.Role{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Description: description,
		Status:      domain.RoleStatusActive,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, tenantID, roleID, name, description string) (*domain.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, tenantID, roleID)
	if err != nil {
		return nil, errors.NewBusinessError(errors.CodeNotFound, "角色不存在")
	}

	role.Name = name
	role.Description = description

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	return role, nil
}

// AssignPermissions 给角色分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, tenantID, roleID string, permissionIDs []string) error {
	if _, err := s.roleRepo.FindByID(ctx, tenantID, roleID); err != nil {
		return errors.NewBusinessError(errors.CodeNotFound, "角色不存在")
	}

	return s.roleRepo.AddPermissions(ctx, roleID, permissionIDs)
}

// GetRole 获取角色详情（含权限列表）
func (s *RoleService) GetRole(ctx context.Context, tenantID, roleID string) (*domain.Role, error) {
	return s.roleRepo.FindWithPermissions(ctx, tenantID, roleID)
}

// ListRoles 分页查询角色列表
func (s *RoleService) ListRoles(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Role, int64, error) {
	return s.roleRepo.List(ctx, tenantID, offset, limit)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, tenantID, roleID string) error {
	if _, err := s.roleRepo.FindByID(ctx, tenantID, roleID); err != nil {
		return errors.NewBusinessError(errors.CodeNotFound, "角色不存在")
	}

	return s.roleRepo.Delete(ctx, tenantID, roleID)
}

// ListPermissions 获取所有权限列表（系统级）
func (s *RoleService) ListPermissions(ctx context.Context, offset, limit int) ([]*domain.Permission, int64, error) {
	return s.permRepo.List(ctx, offset, limit)
}
