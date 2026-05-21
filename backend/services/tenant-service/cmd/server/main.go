package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	cfg.Server.Name = "tenant-service"
	cfg.Server.Port = 8082

	log := logger.New(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Server.Name,
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

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.Server.Name,
		})
	})

	api := engine.Group("/api/v1/tenant")
	{
		api.GET("/tenants", notImpl("租户列表"))
		api.GET("/organizations", notImpl("组织列表"))
		api.GET("/departments", notImpl("部门列表"))
		api.GET("/positions", notImpl("岗位列表"))
	}

	log.Info("Tenant 服务启动（数据库仓储待实现）")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Infof("Tenant 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Tenant 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭 Tenant 服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Tenant 服务关闭异常: %v", err)
	}
	log.Info("Tenant 服务已关闭")
}

func notImpl(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": name + "接口已规划，数据库迁移完成后可用",
		})
	}
}
