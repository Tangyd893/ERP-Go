package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "product-service"
	cfg.Server.Port = 8083

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": cfg.Server.Name})
	})

	api := engine.Group("/api/v1/product")
	{
		api.GET("/spus", notImpl("SPU列表"))
		api.GET("/skus", notImpl("SKU列表"))
		api.GET("/skus/:id", notImpl("SKU详情"))
		api.GET("/sku-mappings", notImpl("SKU映射列表"))
	}

	log.Info("Product 服务启动（数据库仓储待实现）")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * 1000000000, WriteTimeout: 30 * 1000000000}
	go func() {
		log.Infof("Product 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Product 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("正在关闭 Product 服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*1000000000)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info("Product 服务已关闭")
}

func notImpl(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": name + "接口已规划"})
	}
}
