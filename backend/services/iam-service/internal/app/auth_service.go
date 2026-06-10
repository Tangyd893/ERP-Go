package app

import (
	"context"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/google/uuid"
)

const errUserDisabled = "用户已被禁用"

// AuthService 认证应用服务
type AuthService struct {
	userRepo    domain.UserRepository
	roleRepo    domain.RoleRepository
	tokenMgr    domain.TokenManager
	passHasher  domain.PasswordHasher
	auditRepo   domain.AuditRepository
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	tokenMgr domain.TokenManager,
	passHasher domain.PasswordHasher,
	auditRepo domain.AuditRepository,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		tokenMgr:   tokenMgr,
		passHasher: passHasher,
		auditRepo:  auditRepo,
	}
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *domain.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, tenantID, username, password, ip, userAgent string) (*LoginResult, error) {
	user, err := s.userRepo.FindWithRoles(ctx, tenantID, username)
	if err != nil {
		s.writeAudit(ctx, tenantID, "", username, domain.AuditLogin, "user", "", "登录失败", ip, userAgent, "fail", "用户不存在")
		return nil, errors.NewBusinessError(errors.CodeLoginFailed, "用户名或密码错误")
	}

	if !user.IsActive() {
		s.writeAudit(ctx, tenantID, user.ID, username, domain.AuditLogin, "user", user.ID, "登录失败", ip, userAgent, "fail", errUserDisabled)
		return nil, errors.NewBusinessError(errors.CodeUserDisabled, errUserDisabled)
	}

	if !s.passHasher.Verify(password, user.PasswordHash) {
		s.writeAudit(ctx, tenantID, user.ID, username, domain.AuditLogin, "user", user.ID, "登录失败", ip, userAgent, "fail", "密码错误")
		return nil, errors.NewBusinessError(errors.CodeLoginFailed, "用户名或密码错误")
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, tenantID, user.Roles)
	if err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	refreshToken, err := s.tokenMgr.GenerateRefreshToken(user.ID, tenantID)
	if err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	user.RecordLogin()
	_ = s.userRepo.Update(ctx, user)

	s.writeAudit(ctx, tenantID, user.ID, username, domain.AuditLogin, "user", user.ID, "登录成功", ip, userAgent, "success", "")

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    7200,
		User:         user,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	claims, err := s.tokenMgr.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.NewBusinessError(errors.CodeTokenInvalid, "令牌无效或已过期")
	}

	user, err := s.userRepo.FindWithRoles(ctx, claims.TenantID, claims.UserID)
	if err != nil {
		return nil, errors.NewBusinessError(errors.CodeUserDisabled, "用户不存在")
	}

	if !user.IsActive() {
		return nil, errors.NewBusinessError(errors.CodeUserDisabled, errUserDisabled)
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, claims.TenantID, user.Roles)
	if err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	newRefreshToken, err := s.tokenMgr.GenerateRefreshToken(user.ID, claims.TenantID)
	if err != nil {
		return nil, errors.WrapError(errors.CodeInternalError, err)
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    7200,
		User:         user,
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, tenantID, userID, username, ip, userAgent string) error {
	s.writeAudit(ctx, tenantID, userID, username, domain.AuditLogout, "user", userID, "登出", ip, userAgent, "success", "")
	return nil
}

// CheckPermission 检查用户权限
func (s *AuthService) CheckPermission(ctx context.Context, tenantID, userID, permissionCode string) (bool, error) {
	user, err := s.userRepo.FindWithRoles(ctx, tenantID, userID)
	if err != nil {
		return false, errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}

	if !user.IsActive() {
		return false, errors.NewBusinessError(errors.CodeUserDisabled, errUserDisabled)
	}

	return user.HasPermission(permissionCode), nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	user, err := s.userRepo.FindWithRoles(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.NewBusinessError(errors.CodeNotFound, "用户不存在")
	}
	return user, nil
}

func (s *AuthService) writeAudit(ctx context.Context, tenantID, userID, username string, action domain.AuditAction, resourceType, resourceID, detail, ip, userAgent, result, resultDetail string) {
	log := &domain.AuditLog{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detail,
		IP:           ip,
		UserAgent:    userAgent,
		Result:       result,
		ResultDetail: resultDetail,
		CreatedAt:    time.Now(),
	}
	_ = s.auditRepo.Write(ctx, log)
}
