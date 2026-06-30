package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-gateway-secret")
	jwtSecret = []byte("test-gateway-secret")
}

func generateTestToken(userID, tenantID, username string) (string, error) {
	claims := jwtClaims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		Roles:    []string{"admin"},
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "erp-go",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func setupGatewayEngine() *gin.Engine {
	log := logger.New("info", "text", "stdout", "gateway-test", "testing")

	engine := gin.New()
	engine.Use(
		gin.Recovery(),
	)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	engine.Use(authMiddleware(log))

	// 不代理，直接注册一个测试端点验证鉴权
	engine.GET("/api/v1/protected/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")
		c.JSON(http.StatusOK, gin.H{
			"user_id":   userID,
			"tenant_id": tenantID,
		})
	})

	return engine
}

func TestGatewayHealth(t *testing.T) {
	engine := setupGatewayEngine()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Errorf("期望 status=ok，实际 %v", resp["status"])
	}
}

func TestGatewayAuth_MissingToken(t *testing.T) {
	engine := setupGatewayEngine()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestGatewayAuth_InvalidToken(t *testing.T) {
	engine := setupGatewayEngine()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestGatewayAuth_SkipPaths(t *testing.T) {
	// 确认 skipPaths 中的路径依然放行
	engine := setupGatewayEngine()

	// login 路径应跳过鉴权
	req := httptest.NewRequest(http.MethodPost, "/api/v1/iam/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Error("login 路径应跳过鉴权，实际返回 401")
	}

	// health 路径应放行
	req2 := httptest.NewRequest(http.MethodGet, "/health", nil)
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("health 期望 200，实际 %d", w2.Code)
	}
}

func TestGatewayAuth_ValidToken(t *testing.T) {
	// 生成有效 JWT
	token, err := generateTestToken("user-test", "default", "admin")
	if err != nil {
		t.Fatalf("生成测试 token 失败: %v", err)
	}

	engine := setupGatewayEngine()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["user_id"] != "user-test" {
		t.Errorf("期望 user_id=user-test，实际 %v", resp["user_id"])
	}
	if resp["tenant_id"] != "default" {
		t.Errorf("期望 tenant_id=default，实际 %v", resp["tenant_id"])
	}
}

func TestGatewayAuth_MalformedHeader(t *testing.T) {
	engine := setupGatewayEngine()

	// 缺少 Bearer 前缀
	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/test", nil)
	req.Header.Set("Authorization", "some-token-without-bearer")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestGatewayRBAC_AdminAccess(t *testing.T) {
	// 启用 RBAC
	os.Setenv("GATEWAY_RBAC_ENABLED", "true")
	defer os.Unsetenv("GATEWAY_RBAC_ENABLED")

	token, _ := generateTestToken("user-admin", "default", "admin")

	log := logger.New("info", "text", "stdout", "gateway-test", "testing")
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.Use(authMiddleware(log))
	engine.Use(rbacMiddleware(log))
	engine.GET("/api/v1/admin/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"allowed": true})
	})
	engine.GET("/api/v1/order/orders", func(c *gin.Context) {
		c.JSON(200, gin.H{"allowed": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("admin 用户应能访问 /api/v1/admin/，实际 %d", w.Code)
	}
}

func TestGatewayRBAC_NonAdminDenied(t *testing.T) {
	os.Setenv("GATEWAY_RBAC_ENABLED", "true")
	defer os.Unsetenv("GATEWAY_RBAC_ENABLED")

	// 生成只有 "viewer" 角色的 token
	token, _ := generateTestTokenWithRoles("user-viewer", "default", "viewer", []string{"viewer"})

	log := logger.New("info", "text", "stdout", "gateway-test", "testing")
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(authMiddleware(log))
	engine.Use(rbacMiddleware(log))
	engine.GET("/api/v1/admin/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"allowed": true})
	})
	engine.GET("/api/v1/report/reports", func(c *gin.Context) {
		c.JSON(200, gin.H{"allowed": true})
	})

	// viewer 不能访问 admin 路由
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("viewer 访问 admin 路由期望 403，实际 %d", w.Code)
	}

	// viewer 可以访问 report 路由
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/report/reports", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("viewer 访问 report 路由期望 200，实际 %d", w2.Code)
	}
}

func TestGatewayRBAC_SkipPathsUnaffected(t *testing.T) {
	os.Setenv("GATEWAY_RBAC_ENABLED", "true")
	defer os.Unsetenv("GATEWAY_RBAC_ENABLED")

	log := logger.New("info", "text", "stdout", "gateway-test", "testing")
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.Use(authMiddleware(log))
	engine.Use(rbacMiddleware(log))

	// health 不鉴权，RBAC 也放行
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("health 期望 200，实际 %d", w.Code)
	}
}

// generateTestTokenWithRoles 生成指定角色的测试 token
func generateTestTokenWithRoles(userID, tenantID, username string, roles []string) (string, error) {
	claims := jwtClaims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		Roles:    roles,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "erp-go",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
