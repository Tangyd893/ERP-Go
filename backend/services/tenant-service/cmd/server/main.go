package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	cfg.Server.Name = "tenant-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8082
	}

	log := logger.New(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Server.Name,
		os.Getenv("ENVIRONMENT"),
	)

	var db *gorm.DB
	var tenantAppService *app.TenantAppService
	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用占位模式: %v", dbErr)
	} else {
		log.Info("数据库连接成功")
		db = database
		tenantRepo := repository.NewTenantRepository(db)
		orgRepo := repository.NewOrgRepository(db)
		tenantAppService = app.NewTenantAppService(tenantRepo, orgRepo)
	}

	tenantHandler := handler.NewTenantHandler(tenantAppService)

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
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": cfg.Server.Name,
			"db":      db != nil,
		})
	})

	tenantHandler.RegisterRoutes(engine.Group("/api/v1/tenant"))

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

	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Tenant 服务关闭异常: %v", err)
	}
	log.Info("Tenant 服务已关闭")
}
