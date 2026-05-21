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
)

// Route 网关路由配置
type Route struct {
	Path        string `mapstructure:"path"`
	Target      string `mapstructure:"target"`
	StripPrefix bool   `mapstructure:"strip_prefix"`
	AuthRequired bool  `mapstructure:"auth_required"`
}

// GatewayConfig 网关配置
type GatewayConfig struct {
	config.Config
	Routes []Route `mapstructure:"routes"`
}

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

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// 代理路由 - 从配置或环境变量加载
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

func registerProxyRoutes(engine *gin.Engine, log logger.Logger) {
	// 从环境变量加载代理路由，格式：SERVICE_<NAME>_TARGET=<url>, SERVICE_<NAME>_PATH=/api/xxx/
	// 默认路由配置
	defaultRoutes := map[string]string{
		"/api/v1/iam/":   "http://localhost:8081",
		"/api/v1/tenant/": "http://localhost:8082",
		"/api/v1/product/": "http://localhost:8083",
		"/api/v1/channel/": "http://localhost:8084",
		"/api/v1/order/":  "http://localhost:8085",
		"/api/v1/inventory/": "http://localhost:8086",
		"/api/v1/warehouse/": "http://localhost:8087",
		"/api/v1/transport/": "http://localhost:8088",
		"/api/v1/file/":    "http://localhost:8089",
	}

	for path, target := range defaultRoutes {
		// 允许通过环境变量覆盖
		envKey := fmt.Sprintf("SERVICE_TARGET_%s", strings.ToUpper(strings.Trim(path, "/")))
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
			log.Infof("代理请求: %s %s -> %s%s", c.Request.Method, c.Request.URL.Path, targetURL.String(), originalPath)
			proxy.ServeHTTP(c.Writer, c.Request)
		})

		log.Infof("注册代理路由: %s -> %s", proxyPath, target)
	}
}
