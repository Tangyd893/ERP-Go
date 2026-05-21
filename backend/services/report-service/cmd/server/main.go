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

func main() { cfg,_:=config.Load(""); cfg.Server.Name="report-service"; cfg.Server.Port=8093; log:=logger.New(cfg.Log.Level,cfg.Log.Format,cfg.Log.Output,cfg.Server.Name,os.Getenv("ENVIRONMENT"))
	if cfg.Server.Mode=="release" {gin.SetMode(gin.ReleaseMode)}; engine:=gin.New(); engine.Use(middleware.Recovery(log),middleware.RequestID(),middleware.TraceID(),middleware.TenantID(),middleware.CORS(),middleware.RequestLogger(log))
	engine.GET("/health",func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"status":"ok","service":cfg.Server.Name})})
	api:=engine.Group("/api/v1/report")
	api.GET("/sales",notImpl("销售报表"))
	api.GET("/inventory-turnover",notImpl("库存周转"))
	api.GET("/warehouse-efficiency",notImpl("仓储效率"))
	api.GET("/profit-summary",notImpl("利润汇总"))
	log.Info("Report 报表服务启动"); addr:=fmt.Sprintf("%s:%d",cfg.Server.Host,cfg.Server.Port); srv:=&http.Server{Addr:addr,Handler:engine,ReadTimeout:60*time.Second,WriteTimeout:60*time.Second}; go func(){srv.ListenAndServe()}(); select{}}
func notImpl(n string) gin.HandlerFunc { return func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"code":0,"message":n+"接口已规划"})} }
