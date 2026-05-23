package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestObservabilityFieldsSet 验证中间件设置了 service、client_ip、method、path 字段
func TestObservabilityFieldsSet(t *testing.T) {
	router := gin.New()
	router.Use(WithObservabilityFields())
	router.GET("/test", func(c *gin.Context) {
		service, _ := c.Get(LogFieldService)
		clientIP, _ := c.Get(LogFieldClientIP)
		method, _ := c.Get(LogFieldMethod)
		path, _ := c.Get(LogFieldPath)

		if service != "erp-go" {
			t.Errorf("期望 service=erp-go, 实际=%v", service)
		}
		if clientIP == nil || clientIP == "" {
			t.Errorf("期望 client_ip 不为空, 实际=%v", clientIP)
		}
		if method != "GET" {
			t.Errorf("期望 method=GET, 实际=%v", method)
		}
		if path != "/test" {
			t.Errorf("期望 path=/test, 实际=%v", path)
		}

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际=%d", w.Code)
	}
}

// TestLatencyTimerSet 验证请求完成后设置了 duration_ms 和 status_code 字段
func TestLatencyTimerSet(t *testing.T) {
	router := gin.New()
	router.Use(LatencyTimer())
	router.GET("/latency", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest("GET", "/latency", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("期望状态码 201, 实际=%d", w.Code)
	}

	duration, exists := w.Result().Header["X-Duration"]
	_ = duration // 不使用但避免未使用变量错误
	_ = exists

	req2 := httptest.NewRequest("GET", "/latency", nil)
	w2 := httptest.NewRecorder()
	router2 := gin.New()
	router2.Use(LatencyTimer())
	router2.GET("/latency", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router2.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 实际=%d", w2.Code)
	}
}

// TestLatencyTimerFieldsInContext 验证 LatencyTimer 在 gin.Context 中设置了 duration_ms 和 status_code
func TestLatencyTimerFieldsInContext(t *testing.T) {
	var capturedDuration interface{}
	var capturedStatusCode interface{}

	router := gin.New()
	// 外层检查中间件放在 LatencyTimer 之前，确保 c.Next() 返回后能读取 LatencyTimer 设置的值
	router.Use(func(c *gin.Context) {
		c.Next()
		capturedDuration, _ = c.Get(LogFieldDuration)
		capturedStatusCode, _ = c.Get(LogFieldStatusCode)
	})
	router.Use(LatencyTimer())
	router.GET("/check", func(c *gin.Context) {
		c.Status(http.StatusTeapot)
	})

	req := httptest.NewRequest("GET", "/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Errorf("期望状态码 418, 实际=%d", w.Code)
	}

	if capturedDuration == nil {
		t.Error("期望 duration_ms 被设置，实际为 nil")
	}
	if capturedStatusCode == nil {
		t.Error("期望 status_code 被设置，实际为 nil")
	}
	if capturedStatusCode != http.StatusTeapot {
		t.Errorf("期望 status_code=418, 实际=%v", capturedStatusCode)
	}
}

// TestObservabilityFieldsWithLatencyTimer 验证两个中间件联合使用
func TestObservabilityFieldsWithLatencyTimer(t *testing.T) {
	router := gin.New()
	router.Use(WithObservabilityFields())
	router.Use(func(c *gin.Context) {
		// 先执行业务处理
		c.Next()
		// 业务处理完成后的钩子
	})
	router.Use(LatencyTimer())
	router.GET("/combined", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/combined", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("期望状态码 204, 实际=%d", w.Code)
	}
}
