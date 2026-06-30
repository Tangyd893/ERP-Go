package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	httpiface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/gin-gonic/gin"
)

// ── Mocks ──────────────────────────────────────────────

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) Update(_ context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) FindByID(_ context.Context, _, userID string) (*domain.User, error) {
	if u, ok := m.users[userID]; ok {
		return u, nil
	}
	return nil, nil
}
func (m *mockUserRepo) FindByUsername(_ context.Context, _, username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) FindWithRoles(_ context.Context, _, idOrUsername string) (*domain.User, error) {
	// RefreshToken 传入 userID，Login 传入 username，两者都尝试
	u, _ := m.FindByID(nil, "", idOrUsername)
	if u != nil {
		return u, nil
	}
	u, _ = m.FindByUsername(nil, "", idOrUsername)
	return u, nil
}
func (m *mockUserRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.User, int64, error) {
	var list []*domain.User
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, int64(len(list)), nil
}
func (m *mockUserRepo) Delete(_ context.Context, _, userID string) error {
	delete(m.users, userID)
	return nil
}

type mockRoleRepo struct {
	roles       map[string]*domain.Role
	permissions map[string]*domain.Permission
}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{
		roles:       make(map[string]*domain.Role),
		permissions: make(map[string]*domain.Permission),
	}
}

func (m *mockRoleRepo) Create(_ context.Context, role *domain.Role) error {
	m.roles[role.ID] = role
	return nil
}
func (m *mockRoleRepo) Update(_ context.Context, role *domain.Role) error {
	m.roles[role.ID] = role
	return nil
}
func (m *mockRoleRepo) FindByID(_ context.Context, _, roleID string) (*domain.Role, error) {
	if r, ok := m.roles[roleID]; ok {
		return r, nil
	}
	return nil, nil
}
func (m *mockRoleRepo) FindByCode(_ context.Context, _, code string) (*domain.Role, error) {
	for _, r := range m.roles {
		if r.Code == code {
			return r, nil
		}
	}
	return nil, nil
}
func (m *mockRoleRepo) FindWithPermissions(_ context.Context, _, roleID string) (*domain.Role, error) {
	return m.FindByID(nil, "", roleID)
}
func (m *mockRoleRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.Role, int64, error) {
	var list []*domain.Role
	for _, r := range m.roles {
		list = append(list, r)
	}
	return list, int64(len(list)), nil
}
func (m *mockRoleRepo) Delete(_ context.Context, _, roleID string) error {
	delete(m.roles, roleID)
	return nil
}
func (m *mockRoleRepo) AddPermissions(_ context.Context, roleID string, permIDs []string) error {
	return nil
}
func (m *mockRoleRepo) RemovePermissions(_ context.Context, roleID string, permIDs []string) error {
	return nil
}
func (m *mockRoleRepo) AssignUserRoles(_ context.Context, _ string, _ []string) error {
	return nil
}
func (m *mockRoleRepo) RemoveUserRoles(_ context.Context, _ string, _ []string) error {
	return nil
}

type mockPermRepo struct {
	perms map[string]*domain.Permission
}

func newMockPermRepo() *mockPermRepo {
	return &mockPermRepo{perms: make(map[string]*domain.Permission)}
}
func (m *mockPermRepo) Create(_ context.Context, p *domain.Permission) error {
	m.perms[p.ID] = p
	return nil
}
func (m *mockPermRepo) Update(_ context.Context, p *domain.Permission) error {
	m.perms[p.ID] = p
	return nil
}
func (m *mockPermRepo) FindByID(_ context.Context, id string) (*domain.Permission, error) {
	if p, ok := m.perms[id]; ok {
		return p, nil
	}
	return nil, nil
}
func (m *mockPermRepo) FindByCode(_ context.Context, code string) (*domain.Permission, error) {
	for _, p := range m.perms {
		if p.Code == code {
			return p, nil
		}
	}
	return nil, nil
}
func (m *mockPermRepo) List(_ context.Context, _, _ int) ([]*domain.Permission, int64, error) {
	var list []*domain.Permission
	for _, p := range m.perms {
		list = append(list, p)
	}
	return list, int64(len(list)), nil
}
func (m *mockPermRepo) ListByRoleID(_ context.Context, _ string) ([]*domain.Permission, error) {
	return nil, nil
}
func (m *mockPermRepo) Delete(_ context.Context, id string) error {
	delete(m.perms, id)
	return nil
}

type mockAuditRepo struct{}

func (m *mockAuditRepo) Write(_ context.Context, _ *domain.AuditLog) error { return nil }
func (m *mockAuditRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.AuditLog, int64, error) {
	return nil, 0, nil
}

// ── Test Fixture ───────────────────────────────────────

const testJWTSecret = "test-iam-secret-for-httptest"

func setupTestServer(t *testing.T) (*gin.Engine, *mockUserRepo, *mockRoleRepo) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	permRepo := newMockPermRepo()
	auditRepo := &mockAuditRepo{}

	// 先创建角色（用户依赖角色）
	roleRepo.roles["role-1"] = &domain.Role{
		ID:       "role-1",
		TenantID: "default",
		Name:     "超级管理员",
		Code:     "super_admin",
	}

	// 再创建用户（引用已有角色）
	hasher := infra.NewBcryptPasswordHasher(10)
	hash, _ := hasher.Hash("testpass123")
	userRepo.users["user-1"] = &domain.User{
		ID:           "user-1",
		TenantID:     "default",
		Username:     "admin",
		PasswordHash: hash,
		Nickname:     "管理员",
		Status:       "active",
		Roles:        []domain.Role{*roleRepo.roles["role-1"]},
	}

	tokenMgr := infra.NewJWTTokenManager(testJWTSecret, 2*time.Hour, 7*24*time.Hour, "erp-go")

	authSvc := app.NewAuthService(userRepo, roleRepo, tokenMgr, hasher, auditRepo)
	userSvc := app.NewUserService(userRepo, roleRepo, hasher, auditRepo)
	roleSvc := app.NewRoleService(roleRepo, permRepo, auditRepo)

	engine := gin.New()
	log := logger.New("info", "text", "stdout", "iam-test", "testing")
	cfg, _ := config.Load("")

	srv := httpiface.NewServer(authSvc, userSvc, roleSvc, log, cfg, testJWTSecret)
	srv.RegisterRoutes(engine)

	return engine, userRepo, roleRepo
}

func authHeader(token string) string {
	return "Bearer " + token
}

// ── Tests ──────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	body := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 0 {
		t.Fatalf("期望 code=0，实际 %v", resp["code"])
	}
	data := resp["data"].(map[string]interface{})
	if data["token_type"] != "Bearer" {
		t.Errorf("期望 token_type=Bearer，实际 %v", data["token_type"])
	}
	if data["access_token"] == "" {
		t.Error("access_token 不应为空")
	}
	if data["refresh_token"] == "" {
		t.Error("refresh_token 不应为空")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	body := `{"tenant_id":"default","username":"admin","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 200 或 401，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) == 0 {
		t.Fatal("期望非 0 错误码（无效凭证）")
	}
}

func TestLogin_BadRequest(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 缺少必填字段
	body := `{"username":"admin"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望 400，实际 %d", w.Code)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 先登录获取 refresh_token
	loginBody := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	engine.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	refreshToken := loginResp["data"].(map[string]interface{})["refresh_token"].(string)

	// 用 refresh_token 刷新
	body := `{"refresh_token":"` + refreshToken + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 0 {
		t.Fatalf("期望 code=0，实际 %v: %s", resp["code"], w.Body.String())
	}
}

func TestRefreshToken_BadRequest(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("期望 400，实际 %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/user/info", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/user/info", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-xyz")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 登录获取 access_token
	loginBody := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	engine.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	accessToken := loginResp["data"].(map[string]interface{})["access_token"].(string)

	// 用 access_token 访问受保护接口
	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/user/info", nil)
	req.Header.Set("Authorization", authHeader(accessToken))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}
}

func TestPermissionMiddleware_AdminAccess(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 登录获取 access_token
	loginBody := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	engine.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	accessToken := loginResp["data"].(map[string]interface{})["access_token"].(string)

	// super_admin 可以访问 users 列表
	req := httptest.NewRequest(http.MethodGet, "/api/v1/iam/users", nil)
	req.Header.Set("Authorization", authHeader(accessToken))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}
}

func TestLogout_Success(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 登录
	loginBody := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	engine.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	accessToken := loginResp["data"].(map[string]interface{})["access_token"].(string)

	// 登出
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/logout", nil)
	req.Header.Set("Authorization", authHeader(accessToken))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}
}

func TestCheckPermission_Success(t *testing.T) {
	engine, _, _ := setupTestServer(t)

	// 登录（admin 已有 super_admin 角色）
	loginBody := `{"tenant_id":"default","username":"admin","password":"testpass123"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	engine.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	accessToken := loginResp["data"].(map[string]interface{})["access_token"].(string)

	// 检查权限
	body := `{"user_id":"user-1","tenant_id":"default","permission_code":"order:read"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/check-permission", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader(accessToken))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}
}
