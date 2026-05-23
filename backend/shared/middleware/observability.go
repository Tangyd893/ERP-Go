package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// 统一日志字段名
	LogFieldService    = "service"
	LogFieldTraceID    = "trace_id"
	LogFieldRequestID  = "request_id"
	LogFieldTenantID   = "tenant_id"
	LogFieldUserID     = "user_id"
	LogFieldBusinessNo = "business_no"
	LogFieldDuration   = "duration_ms"
	LogFieldStatusCode = "status_code"
	LogFieldMethod     = "method"
	LogFieldPath       = "path"
	LogFieldClientIP   = "client_ip"
)

// WithObservabilityFields 将统一可观测字段注入请求上下文
func WithObservabilityFields() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(LogFieldService, "erp-go")
		c.Set(LogFieldClientIP, c.ClientIP())
		c.Set(LogFieldMethod, c.Request.Method)
		c.Set(LogFieldPath, c.Request.URL.Path)
		c.Next()
	}
}

// LatencyTimer 请求延迟计时（在中间件链末尾使用）
func LatencyTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Milliseconds()
		c.Set(LogFieldDuration, duration)
		c.Set(LogFieldStatusCode, c.Writer.Status())
	}
}

// GetObservabilityField 从上下文获取可观测字段值
func GetObservabilityField(c *gin.Context, key string) interface{} {
	if val, exists := c.Get(key); exists {
		return val
	}
	return nil
}

// ObservabilityContextKey 可观测上下文键类型
type ObservabilityContextKey struct{}

// InjectObservabilityToContext 将可观测字段注入标准 context.Context
func InjectObservabilityToContext(ctx context.Context, c *gin.Context) context.Context {
	for _, key := range []string{
		LogFieldService,
		LogFieldTraceID,
		LogFieldRequestID,
		LogFieldTenantID,
		LogFieldUserID,
		LogFieldBusinessNo,
		LogFieldClientIP,
		LogFieldMethod,
		LogFieldPath,
	} {
		if val, exists := c.Get(key); exists {
			ctx = context.WithValue(ctx, ObservabilityContextKey{}, map[string]interface{}{
				key: val,
			})
		}
	}
	return ctx
}
