package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/infra/repository"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	mu           sync.RWMutex
	mockBalances map[string]*domain.InventoryBalance
	inventoryRepo *repository.InventoryRepository
	db           *gorm.DB
)

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "inventory-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8086
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用内存模拟模式: %v", dbErr)
		mockBalances = map[string]*domain.InventoryBalance{
			"sku-001": {ID: "bal-001", TenantID: "t-001", WarehouseID: "wh-001", SKUID: "sku-001", SKUCode: "TSHIRT-001", TotalQuantity: 500, Version: 1},
			"sku-002": {ID: "bal-002", TenantID: "t-001", WarehouseID: "wh-001", SKUID: "sku-002", SKUCode: "MUG-001", TotalQuantity: 300, Version: 1},
		}
	} else {
		log.Info("数据库连接成功")
		db = database
		inventoryRepo = repository.NewInventoryRepository(db)
	}

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

	api := engine.Group("/api/v1/inventory")
	{
		api.GET("/balances", listBalances)
		api.GET("/balances/:sku_id", getBalance)
		api.POST("/lock", lockInventory)
		api.POST("/release", releaseInventory)
		api.POST("/deduct", deductInventory)
		api.GET("/journals", listJournals)
	}

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

func listBalances(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()
	balances := make([]*domain.InventoryBalance, 0, len(mockBalances))
	for _, v := range mockBalances {
		balances = append(balances, v)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": balances})
}

func getBalance(c *gin.Context) {
	skuID := c.Param("sku_id")
	mu.RLock()
	bal, ok := mockBalances[skuID]
	mu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 40000, "message": "SKU库存未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": bal})
}

func lockInventory(c *gin.Context) {
	var req struct {
		SKUID   string `json:"sku_id" binding:"required"`
		OrderID string `json:"order_id" binding:"required"`
		Qty     int    `json:"quantity" binding:"required"`
		LockKey string `json:"lock_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 10001, "message": "参数无效"})
		return
	}

	mu.Lock()
	bal, ok := mockBalances[req.SKUID]
	if !ok {
		mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"code": 40000, "message": "SKU库存未找到"})
		return
	}
	if err := bal.Lock(req.Qty); err != nil {
		mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 60001, "message": err.Error()})
		return
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"locked":          true,
			"available":       bal.Available(),
			"locked_quantity": bal.LockedQuantity,
		},
	})
}

func releaseInventory(c *gin.Context) {
	var req struct {
		SKUID string `json:"sku_id" binding:"required"`
		Qty   int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 10001, "message": "参数无效"})
		return
	}

	mu.Lock()
	bal, ok := mockBalances[req.SKUID]
	if !ok {
		mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"code": 40000, "message": "SKU库存未找到"})
		return
	}
	if err := bal.Release(req.Qty); err != nil {
		mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 60002, "message": err.Error()})
		return
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"released": true, "available": bal.Available()},
	})
}

func deductInventory(c *gin.Context) {
	var req struct {
		SKUID string `json:"sku_id" binding:"required"`
		Qty   int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 10001, "message": "参数无效"})
		return
	}

	mu.Lock()
	bal, ok := mockBalances[req.SKUID]
	if !ok {
		mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"code": 40000, "message": "SKU库存未找到"})
		return
	}
	if err := bal.Deduct(req.Qty); err != nil {
		mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 60003, "message": err.Error()})
		return
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"deducted": true, "total": bal.TotalQuantity},
	})
}

func listJournals(c *gin.Context) {
	skuID := c.Query("sku_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	journals := []*domain.InventoryJournal{
		{ID: "j-1", SKUID: coa(skuID, "sku-001"), ChangeType: "lock", ChangeQty: 10, BeforeTotal: 500, AfterTotal: 500, BeforeLocked: 0, AfterLocked: 10, BeforeAvail: 500, AfterAvail: 490},
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": journals, "total": len(journals), "page": page, "page_size": pageSize}})
}

func coa(val, defaultVal string) string {
	if strings.TrimSpace(val) == "" {
		return defaultVal
	}
	return val
}
