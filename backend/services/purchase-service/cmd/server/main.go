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

func main() { cfg,_:=config.Load(""); cfg.Server.Name="purchase-service"; cfg.Server.Port=8091; log:=logger.New(cfg.Log.Level,cfg.Log.Format,cfg.Log.Output,cfg.Server.Name,os.Getenv("ENVIRONMENT"))
	if cfg.Server.Mode=="release" {gin.SetMode(gin.ReleaseMode)}; engine:=gin.New(); engine.Use(middleware.Recovery(log),middleware.RequestID(),middleware.TraceID(),middleware.TenantID(),middleware.CORS(),middleware.RequestLogger(log))
	engine.GET("/health",func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"status":"ok","service":cfg.Server.Name})})
	api:=engine.Group("/api/v1/purchase")
	api.GET("/suppliers",notImpl("供应商列表"))
	api.POST("/suppliers",notImpl("添加供应商"))
	api.GET("/orders",notImpl("采购单列表"))
	api.POST("/orders",notImpl("创建采购单"))
	api.GET("/inbound",notImpl("入库单列表"))
	log.Info("Purchase 采购服务启动"); addr:=fmt.Sprintf("%s:%d",cfg.Server.Host,cfg.Server.Port); srv:=&http.Server{Addr:addr,Handler:engine,ReadTimeout:30*time.Second,WriteTimeout:30*time.Second}; go func(){srv.ListenAndServe()}(); select{}}
func notImpl(n string) gin.HandlerFunc { return func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"code":0,"message":n+"接口已规划"})} }
