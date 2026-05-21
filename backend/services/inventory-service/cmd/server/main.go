package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
)

// 模拟内存存储（生产环境使用数据库）
var (
	mockBalances = map[string]*domain.InventoryBalance{}
	mu           sync.RWMutex
)

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "inventory-service"
	cfg.Server.Port = 8086

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	// 初始化模拟库存数据
	mockBalances["sku-001"] = &domain.InventoryBalance{
		ID: "bal-001", TenantID: "t-001", WarehouseID: "wh-001",
		SKUID: "sku-001", SKUCode: "TSHIRT-001", TotalQuantity: 500, Version: 1,
	}
	mockBalances["sku-002"] = &domain.InventoryBalance{
		ID: "bal-002", TenantID: "t-001", WarehouseID: "wh-001",
		SKUID: "sku-002", SKUCode: "MUG-001", TotalQuantity: 300, Version: 1,
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": cfg.Server.Name})
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

	log.Info("Inventory 服务启动")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}
	go func() { log.Infof("Inventory 服务启动在 %s", addr); srv.ListenAndServe() }()

	select {}
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
		SKUID       string `json:"sku_id" binding:"required"`
		OrderID     string `json:"order_id" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
		LockKey     string `json:"lock_key" binding:"required"`
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
	if err := bal.Lock(req.Quantity); err != nil {
		mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"code": 60001, "message": err.Error()})
		return
	}
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"locked": true,
			"available": bal.Available(),
			"locked_quantity": bal.LockedQuantity,
		},
	})
}

func releaseInventory(c *gin.Context) {
	var req struct {
		SKUID    string `json:"sku_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
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
	if err := bal.Release(req.Quantity); err != nil {
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
		SKUID    string `json:"sku_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
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
	if err := bal.Deduct(req.Quantity); err != nil {
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

	// 模拟数据
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
