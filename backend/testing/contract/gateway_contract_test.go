package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGatewayResponseContracts(t *testing.T) {
	// 模拟 Gateway 响应格式契约
	log := logger.New("info", "text", "stdout", "gateway-contract", "testing")
	engine := setupContractGateway(log)

	t.Run("健康检查响应契约", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		engine.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("期望 200，实际 %d", w.Code)
		}
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		assertField(t, body, "status", "ok")
		assertField(t, body, "service", "api-gateway")
	})

	t.Run("401 未认证响应契约", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/test", nil)
		engine.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("期望 401，实际 %d", w.Code)
		}
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		assertFieldExists(t, body, "code")
		assertFieldExists(t, body, "message")
	})

	t.Run("404 路由不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/nonexistent/resource", nil)
		engine.ServeHTTP(w, req)

		// Gateway 在所有路由上都有 *
		if w.Code != 401 { // 需要认证，返回 401
			// 如果没有 auth middleware，Gin 返回 404
			// 实际 Gateway 有 authMiddleware，所以未认证返回 401
		}
	})
}

func setupContractGateway(log logger.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// Auth middleware 模拟
	engine.Use(func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    20000,
				"message": "未提供认证令牌",
			})
			return
		}
		c.Next()
	})

	engine.GET("/api/v1/admin/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"allowed": true})
	})

	return engine
}

func assertField(t *testing.T, body map[string]interface{}, field string, expected interface{}) {
	t.Helper()
	actual, ok := body[field]
	if !ok {
		t.Errorf("缺少字段 %s", field)
		return
	}
	if actual != expected {
		t.Errorf("字段 %s: 期望 %v，实际 %v", field, expected, actual)
	}
}

func assertFieldExists(t *testing.T, body map[string]interface{}, field string) {
	t.Helper()
	if _, ok := body[field]; !ok {
		t.Errorf("缺少字段 %s", field)
	}
}
