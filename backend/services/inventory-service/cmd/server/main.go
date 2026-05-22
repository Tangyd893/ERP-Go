package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "inventory-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8086
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	var inventoryRepo *repository.InventoryRepository
	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用内存模拟模式: %v", dbErr)
	} else {
		log.Info("数据库连接成功")
		db = database
		inventoryRepo = repository.NewInventoryRepository(db)
	}

	inventoryHandler := handler.NewInventoryHandler(inventoryRepo)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{"status": status, "service": cfg.Server.Name, "db": db != nil})
	})

	inventoryHandler.RegisterRoutes(engine.Group("/api/v1/inventory"))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}

	go func() {
		log.Infof("Inventory 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Inventory 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭 Inventory 服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Inventory 服务关闭异常: %v", err)
	}
	log.Info("Inventory 服务已关闭")
}
