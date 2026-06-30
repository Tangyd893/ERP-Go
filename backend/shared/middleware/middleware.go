package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 请求 ID 中间件，每个请求分配唯一请求 ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// TraceID 追踪 ID 中间件，用于链路追踪
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}

// TenantID 租户 ID 中间件，从请求头提取租户信息
func TenantID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID != "" {
			c.Set("tenant_id", tenantID)
		}
		c.Next()
	}
}

// UserID 用户 ID 中间件，从请求上下文提取用户信息（鉴权后设置）
func UserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()

		fields := logger.Fields{
			"method":      c.Request.Method,
			"path":        path,
			"query":       query,
			"status":      statusCode,
			"duration_ms": duration,
			"client_ip":   c.ClientIP(),
		}

		if requestID, exists := c.Get("request_id"); exists {
			fields["request_id"] = requestID
		}
		if traceID, exists := c.Get("trace_id"); exists {
			fields["trace_id"] = traceID
		}
		if tenantID, exists := c.Get("tenant_id"); exists {
			fields["tenant_id"] = tenantID
		}
		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}

		l := log.WithFields(fields).
			WithField(logger.FieldDuration, duration)

		if statusCode >= 500 {
			l.Error("请求处理异常")
		} else if statusCode >= 400 {
			l.Warn("请求处理警告")
		} else {
			l.Info("请求处理完成")
		}
	}
}

// Recovery 恐慌恢复中间件
func Recovery(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(logger.Fields{
					"panic": err,
					"path":  c.Request.URL.Path,
				}).Error("服务发生恐慌")
				c.AbortWithStatusJSON(500, gin.H{
					"code":    10000,
					"message": "系统内部错误",
				})
			}
		}()
		c.Next()
	}
}

// CORS 跨域中间件
// 开发环境默认允许所有来源；生产环境通过 CORS_ALLOWED_ORIGINS 环境变量配置白名单（逗号分隔）
func CORS() gin.HandlerFunc {
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	allowAll := allowedOrigins == "*"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			// 校验请求 Origin 是否在白名单中
			for _, o := range strings.Split(allowedOrigins, ",") {
				if strings.TrimSpace(o) == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Vary", "Origin")
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Trace-ID,X-Tenant-ID,X-User-ID,X-Idempotency-Key")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID,X-Trace-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GetRequestID 从上下文获取请求 ID
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return ""
}

// GetTraceID 从上下文获取追踪 ID
func GetTraceID(c *gin.Context) string {
	if id, exists := c.Get("trace_id"); exists {
		return id.(string)
	}
	return ""
}

// GetTenantID 从上下文获取租户 ID
func GetTenantID(c *gin.Context) string {
	if id, exists := c.Get("tenant_id"); exists {
		return id.(string)
	}
	return ""
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) string {
	if id, exists := c.Get("user_id"); exists {
		return id.(string)
	}
	return ""
}
