package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/google/uuid"
)

const errUserDisabled = "用户已被禁用"

const sessionTTL = 2 * time.Hour

// SessionStore 可选会话缓存（Redis 试点，T-603）。
type SessionStore interface {
	Put(ctx context.Context, sessionID, userJSON string, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
}

// AuthService 认证应用服务
type AuthService struct {
	userRepo   domain.UserRepository
	roleRepo   domain.RoleRepository
	tokenMgr   domain.TokenManager
	passHasher domain.PasswordHasher
	auditRepo  domain.AuditRepository
	sessions   SessionStore
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

// WithSessionStore 注入 Redis 会话缓存；未注入时登录/登出行为不变。
func (s *AuthService) WithSessionStore(store SessionStore) *AuthService {
	s.sessions = store
	return s
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *domain.User `json:"user"`
}

type sessionPayload struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Username string `json:"username"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, tenantID, username, password, ip, userAgent string) (*LoginResult, error) {
	user, err := s.userRepo.FindWithRoles(ctx, tenantID, username)
	if err != nil {
		s.writeAudit(ctx, &domain.AuditLog{
			TenantID: tenantID, Username: username,
			Action: domain.AuditLogin, ResourceType: "user",
			Detail: "登录失败", IP: ip, UserAgent: userAgent,
			Result: "fail", ResultDetail: "用户不存在",
		})
		return nil, errors.NewBusinessError(errors.CodeLoginFailed, "用户名或密码错误")
	}

	if !user.IsActive() {
		s.writeAudit(ctx, &domain.AuditLog{
			TenantID: tenantID, UserID: user.ID, Username: username,
			Action: domain.AuditLogin, ResourceType: "user", ResourceID: user.ID,
			Detail: "登录失败", IP: ip, UserAgent: userAgent,
			Result: "fail", ResultDetail: errUserDisabled,
		})
		return nil, errors.NewBusinessError(errors.CodeUserDisabled, errUserDisabled)
	}

	if !s.passHasher.Verify(password, user.PasswordHash) {
		s.writeAudit(ctx, &domain.AuditLog{
			TenantID: tenantID, UserID: user.ID, Username: username,
			Action: domain.AuditLogin, ResourceType: "user", ResourceID: user.ID,
			Detail: "登录失败", IP: ip, UserAgent: userAgent,
			Result: "fail", ResultDetail: "密码错误",
		})
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

	s.putSession(ctx, user)

	s.writeAudit(ctx, &domain.AuditLog{
		TenantID: tenantID, UserID: user.ID, Username: username,
		Action: domain.AuditLogin, ResourceType: "user", ResourceID: user.ID,
		Detail: "登录成功", IP: ip, UserAgent: userAgent,
		Result: "success",
	})

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

	s.putSession(ctx, user)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    7200,
		User:         user,
	}, nil
}

// Logout 登出（可选吊销 refresh_token，并清除会话缓存）
func (s *AuthService) Logout(ctx context.Context, tenantID, userID, username, ip, userAgent, refreshToken string) error {
	if refreshToken != "" {
		_ = s.tokenMgr.RevokeRefreshToken(refreshToken)
	}
	if s.sessions != nil && userID != "" {
		_ = s.sessions.Delete(ctx, userID)
	}
	s.writeAudit(ctx, &domain.AuditLog{
		TenantID: tenantID, UserID: userID, Username: username,
		Action: domain.AuditLogout, ResourceType: "user", ResourceID: userID,
		Detail: "登出", IP: ip, UserAgent: userAgent,
		Result: "success",
	})
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

func (s *AuthService) putSession(ctx context.Context, user *domain.User) {
	if s.sessions == nil || user == nil {
		return
	}
	body, err := json.Marshal(sessionPayload{
		UserID: user.ID, TenantID: user.TenantID, Username: user.Username,
	})
	if err != nil {
		return
	}
	_ = s.sessions.Put(ctx, user.ID, string(body), sessionTTL)
}

func (s *AuthService) writeAudit(ctx context.Context, log *domain.AuditLog) {
	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()
	_ = s.auditRepo.Write(ctx, log)
}
