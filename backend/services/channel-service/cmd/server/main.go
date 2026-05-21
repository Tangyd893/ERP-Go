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
	cfg.Server.Name = "channel-service"
	cfg.Server.Port = 8084

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": cfg.Server.Name})
	})

	api := engine.Group("/api/v1/channel")
	{
		api.GET("/stores", notImpl("店铺列表"))
		api.POST("/stores", notImpl("添加店铺"))
		api.POST("/orders/import", notImpl("订单导入"))
		api.GET("/import-tasks", notImpl("导入任务列表"))
		api.GET("/sync-tasks", notImpl("同步任务列表"))
	}

	log.Info("Channel 服务启动（数据库仓储待实现）")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * 1000000000, WriteTimeout: 30 * 1000000000}
	go func() { srv.ListenAndServe() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Channel 服务关闭中...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*1000000000)
	defer cancel()
	srv.Shutdown(ctx)
}

func notImpl(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": name + "接口已规划"})
	}
}
