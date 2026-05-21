package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "notification-service"
	cfg.Server.Port = 8094
	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	if cfg.Server.Mode == "release" { gin.SetMode(gin.ReleaseMode) }
	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": cfg.Server.Name}) })

	api := engine.Group("/api/v1/notification")
	{
		api.GET("/list", notImpl("通知列表"))
		api.PUT("/read", notImpl("标记已读"))
		api.GET("/unread-count", notImpl("未读数"))
	}

	log.Info("Notification 通知服务启动")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}
	go func() { srv.ListenAndServe() }()
	select {}
}
func notImpl(name string) gin.HandlerFunc { return func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"code": 0, "message": name + "接口已规划"}) } }
