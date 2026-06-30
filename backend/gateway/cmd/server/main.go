package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	log := logger.New(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		"api-gateway",
		os.Getenv("ENVIRONMENT"),
	)

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(config.DefaultJWTSecret)
		log.Warn("使用默认 JWT Secret，生产环境请设置 JWT_SECRET 环境变量")
	}

	// T-609: 生产环境拒绝默认 JWT Secret
	if config.IsProduction() {
		if err := config.ValidateProduction(string(jwtSecret), cfg.Database.Password); err != nil {
			log.Fatalf("生产环境安全校验失败: %v", err)
		}
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(
		middleware.Recovery(log),
		middleware.RequestID(),
		middleware.TraceID(),
		middleware.TenantID(),
		middleware.UserID(),
		middleware.CORS(),
		middleware.RequestLogger(log),
	)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	engine.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!DOCTYPE html>
<html lang="zh-CN"><head><meta charset="UTF-8"><title>ERP-Go API Gateway</title>
<style>body{font-family:sans-serif;max-width:640px;margin:40px auto;padding:0 16px;line-height:1.6}
a{color:#409EFF}code{background:#f5f7fa;padding:2px 6px;border-radius:4px}</style></head>
<body>
<h1>ERP-Go API Gateway</h1>
<p>8080 为<strong>后端 API 入口</strong>，不提供 Web 界面。请访问前端开发服务器：</p>
<ul>
<li><a href="http://localhost:5173/">Admin 管理后台</a>（5173）</li>
<li><a href="http://localhost:5174/">Warehouse PDA</a>（5174）</li>
<li><a href="http://localhost:5175/">Dashboard 看板</a>（5175）</li>
</ul>
<p>健康检查：<a href="/health"><code>/health</code></a></p>
<p>登录接口：<code>POST /api/v1/iam/login</code>（admin / admin123，tenant=default）</p>
</body></html>`)
	})

	engine.Use(authMiddleware(log))

	// T-608: 可选 RBAC 中间件（GATEWAY_RBAC_ENABLED=true 时启用）
	if os.Getenv("GATEWAY_RBAC_ENABLED") == "true" {
		engine.Use(rbacMiddleware(log))
	}

	registerProxyRoutes(engine, log)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Infof("API 网关启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("网关启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭网关...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("网关关闭异常: %v", err)
	}
	log.Info("网关已关闭")
}

func authMiddleware(log logger.Logger) gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/":                  true,
		"/health":            true,
		"/api/v1/iam/login":  true,
		"/api/v1/iam/refresh": true,
	}

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    20000,
				"message": "未提供认证令牌",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseJWT(tokenString)
		if err != nil {
			log.Warnf("JWT 验证失败: %v, path=%s", err, c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    20002,
				"message": "令牌无效或已过期",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Tenant-ID", claims.TenantID)
		c.Request.Header.Set("X-Username", claims.Username)

		c.Next()
	}
}

type jwtClaims struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"`
	jwt.RegisteredClaims
}

func parseJWT(tokenString string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的令牌")
	}
	return claims, nil
}

// rbacMiddleware 基于 JWT Roles 的粗粒度 RBAC（T-608）
// 通过 GATEWAY_RBAC_ENABLED=true 启用
// 路由权限映射可通过 GATEWAY_RBAC_RULES 环境变量配置（JSON 格式）
func rbacMiddleware(log logger.Logger) gin.HandlerFunc {
	// 默认权限映射：路径前缀 → 需要的角色
	defaultRules := map[string]string{
		"/api/v1/admin/":       "admin",
		"/api/v1/finance/":     "finance",
		"/api/v1/report/":      "viewer",
		"/api/v1/notification/": "admin",
	}

	// 环境变量可覆盖
	var rules map[string]string
	if raw := os.Getenv("GATEWAY_RBAC_RULES"); raw != "" {
		// 简单格式: "path1=role1,path2=role2"
		rules = make(map[string]string)
		for _, pair := range strings.Split(raw, ",") {
			kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(kv) == 2 {
				rules[kv[0]] = kv[1]
			}
		}
	} else {
		rules = defaultRules
	}

	// 完全放行的路径
	skipPaths := map[string]bool{
		"/health":                   true,
		"/api/v1/iam/login":         true,
		"/api/v1/iam/refresh":       true,
		"/api/v1/iam/user/info":     true,
		"/api/v1/iam/logout":        true,
		"/api/v1/iam/check-permission": true,
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if skipPaths[path] {
			c.Next()
			return
		}

		rolesVal, exists := c.Get("roles")
		if !exists {
			// 尝试从 claims 中获取
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    20003,
				"message": "RBAC 未获取到角色信息",
			})
			return
		}

		userRoles, ok := rolesVal.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    20003,
				"message": "角色信息格式错误",
			})
			return
		}

		// super_admin 完全放行
		for _, r := range userRoles {
			if r == "super_admin" || r == "admin" {
				c.Next()
				return
			}
		}

		// 检查路由权限
		for prefix, requiredRole := range rules {
			if strings.HasPrefix(path, prefix) {
				allowed := false
				for _, r := range userRoles {
					if r == requiredRole {
						allowed = true
						break
					}
				}
				if !allowed {
					log.Warnf("RBAC 拒绝: user=%v path=%s need=%s have=%v",
						c.GetString("username"), path, requiredRole, userRoles)
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code":    20003,
						"message": fmt.Sprintf("权限不足: 需要角色 %s", requiredRole),
					})
					return
				}
			}
		}

		c.Next()
	}
}

func registerProxyRoutes(engine *gin.Engine, log logger.Logger) {
	defaultRoutes := map[string]string{
		"/api/v1/iam/":          "http://localhost:8081",
		"/api/v1/tenant/":       "http://localhost:8082",
		"/api/v1/product/":      "http://localhost:8083",
		"/api/v1/channel/":      "http://localhost:8084",
		"/api/v1/order/":        "http://localhost:8085",
		"/api/v1/inventory/":    "http://localhost:8086",
		"/api/v1/warehouse/":    "http://localhost:8087",
		"/api/v1/transport/":    "http://localhost:8088",
		"/api/v1/file/":         "http://localhost:8089",
		"/api/v1/purchase/":     "http://localhost:8091",
		"/api/v1/finance/":      "http://localhost:8092",
		"/api/v1/report/":       "http://localhost:8093",
		"/api/v1/notification/": "http://localhost:8094",
	}

	for path, target := range defaultRoutes {
		envKey := fmt.Sprintf("SERVICE_TARGET_%s", strings.ToUpper(strings.ReplaceAll(strings.Trim(path, "/"), "/", "_")))
		if envTarget := os.Getenv(envKey); envTarget != "" {
			target = envTarget
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			log.Warnf("无效的代理目标 %s: %v", target, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxyPath := path

		engine.Any(proxyPath+"*path", func(c *gin.Context) {
			originalPath := c.Param("path")
			log.Debugf("代理请求: %s %s -> %s%s", c.Request.Method, c.Request.URL.Path, targetURL.String(), originalPath)
			proxy.ServeHTTP(c.Writer, c.Request)
		})

		log.Infof("注册代理路由: %s -> %s", proxyPath, target)
	}
}
