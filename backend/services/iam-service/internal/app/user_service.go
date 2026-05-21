package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/google/uuid"
)

// UserService 用户管理应用服务
type UserService struct {
	userRepo   domain.UserRepository
	roleRepo   domain.RoleRepository
	passHasher domain.PasswordHasher
	auditRepo  domain.AuditRepository
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	passHasher domain.PasswordHasher,
	auditRepo domain.AuditRepository,
) *UserService {
	return &UserService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		passHasher: passHasher,
		auditRepo:  auditRepo,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, tenantID, username, password, nickname, email, phone string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByUsername(ctx, tenantID, username)
	if existing != nil {
		return nil, errors.NewBusinessError(errors.CodeAlreadyExists, "用户名已存在")
	}

	hash, err := s.passHasher.Hash(password)
	if err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		Username:     username,
		PasswordHash: hash,
		Nickname:     nickname,
		Email:        email,
		Phone:        phone,
		Status:       domain.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, tenantID, userID, nickname, email, phone string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}

	user.Nickname = nickname
	user.Email = email
	user.Phone = phone

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	return user, nil
}

// AssignRoles 给用户分配角色
func (s *UserService) AssignRoles(ctx context.Context, tenantID, userID string, roleIDs []string) error {
	if _, err := s.userRepo.FindByID(ctx, tenantID, userID); err != nil {
		return errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}

	return s.roleRepo.AssignUserRoles(ctx, userID, roleIDs)
}

// DisableUser 禁用用户
func (s *UserService) DisableUser(ctx context.Context, tenantID, userID string) error {
	user, err := s.userRepo.FindByID(ctx, tenantID, userID)
	if err != nil {
		return errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}

	user.Disable()
	return s.userRepo.Update(ctx, user)
}

// EnableUser 启用用户
func (s *UserService) EnableUser(ctx context.Context, tenantID, userID string) error {
	user, err := s.userRepo.FindByID(ctx, tenantID, userID)
	if err != nil {
		return errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}

	user.Enable()
	return s.userRepo.Update(ctx, user)
}

// ListUsers 分页查询用户列表
func (s *UserService) ListUsers(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, int64, error) {
	return s.userRepo.List(ctx, tenantID, offset, limit)
}

// GetUser 获取用户详情
func (s *UserService) GetUser(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	return s.userRepo.FindWithRoles(ctx, tenantID, userID)
}
