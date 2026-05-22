package http

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/infra/repository"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

// InventoryHandler 库存 HTTP 处理器
type InventoryHandler struct {
	repo        *repository.InventoryRepository
	mu          sync.RWMutex
	mockEnabled bool
	mockBalances map[string]*domain.InventoryBalance
}

// NewInventoryHandler 创建库存处理器，repo 为 nil 时自动使用内存模拟模式
func NewInventoryHandler(repo *repository.InventoryRepository) *InventoryHandler {
	h := &InventoryHandler{repo: repo}
	if repo == nil {
		h.mockEnabled = true
		h.mockBalances = map[string]*domain.InventoryBalance{
			"sku-001": {ID: "bal-001", TenantID: "t-001", WarehouseID: "wh-001", SKUID: "sku-001", SKUCode: "TSHIRT-001", TotalQuantity: 500, Version: 1},
			"sku-002": {ID: "bal-002", TenantID: "t-001", WarehouseID: "wh-001", SKUID: "sku-002", SKUCode: "MUG-001", TotalQuantity: 300, Version: 1},
		}
	}
	return h
}

// RegisterRoutes 注册路由
func (h *InventoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/balances", h.listBalances)
	router.GET("/balances/:sku_id", h.getBalance)
	router.POST("/lock", h.lockInventory)
	router.POST("/release", h.releaseInventory)
	router.POST("/deduct", h.deductInventory)
	router.GET("/journals", h.listJournals)
}

func (h *InventoryHandler) listBalances(c *gin.Context) {
	if h.mockEnabled {
		h.listBalancesMock(c)
		return
	}

	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	balances, total, err := h.repo.ListBalances(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询库存失败: "+err.Error())
		return
	}
	response.PageSuccess(c, balances, total, page, pageSize)
}

func (h *InventoryHandler) listBalancesMock(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	balances := make([]*domain.InventoryBalance, 0, len(h.mockBalances))
	for _, v := range h.mockBalances {
		balances = append(balances, v)
	}
	response.Success(c, balances)
}

func (h *InventoryHandler) getBalance(c *gin.Context) {
	if h.mockEnabled {
		h.getBalanceMock(c)
		return
	}
	skuID := c.Param("sku_id")
	warehouseID := c.DefaultQuery("warehouse_id", "wh-001")

	bal, err := h.repo.FindBalance(c.Request.Context(), warehouseID, skuID)
	if err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, bal)
}

func (h *InventoryHandler) getBalanceMock(c *gin.Context) {
	skuID := c.Param("sku_id")
	h.mu.RLock()
	bal, ok := h.mockBalances[skuID]
	h.mu.RUnlock()
	if !ok {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeSKUNotFound, "SKU库存未找到")
		return
	}
	response.Success(c, bal)
}

func (h *InventoryHandler) lockInventory(c *gin.Context) {
	var req struct {
		SKUID       string `json:"sku_id" binding:"required"`
		OrderID     string `json:"order_id" binding:"required"`
		WarehouseID string `json:"warehouse_id"`
		Qty         int    `json:"quantity" binding:"required"`
		LockKey     string `json:"lock_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	if h.mockEnabled {
		h.lockInventoryMock(c, req.SKUID, req.OrderID, req.Qty, req.LockKey)
		return
	}

	tenantID := c.GetString("tenant_id")
	warehouseID := req.WarehouseID
	if warehouseID == "" {
		warehouseID = "wh-001"
	}

	lock := &domain.InventoryLock{
		ID:          req.LockKey + "-" + req.SKUID,
		TenantID:    tenantID,
		OrderID:     req.OrderID,
		SKUID:       req.SKUID,
		WarehouseID: warehouseID,
		Quantity:    req.Qty,
		Status:      "locked",
		LockKey:     req.LockKey,
		CreatedAt:   time.Now(),
	}

	journal := &domain.InventoryJournal{
		ID:             "jrnl-" + req.LockKey + "-" + req.SKUID,
		TenantID:       tenantID,
		WarehouseID:    warehouseID,
		SKUID:          req.SKUID,
		OrderID:        req.OrderID,
		ChangeType:     "lock",
		ChangeQty:      req.Qty,
		IdempotencyKey: req.LockKey,
		Operator:       c.GetString("username"),
		CreatedAt:      time.Now(),
	}

	ctx := c.Request.Context()
	if err := h.repo.LockStock(ctx, lock, journal); err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}

	bal, _ := h.repo.FindBalance(ctx, warehouseID, req.SKUID)
	available := 0
	if bal != nil {
		available = bal.Available()
	}

	response.Success(c, gin.H{
		"locked":          true,
		"available":       available,
		"locked_quantity": req.Qty,
	})
}

func (h *InventoryHandler) lockInventoryMock(c *gin.Context, skuID, orderID string, qty int, lockKey string) {
	_ = orderID
	_ = lockKey

	h.mu.Lock()
	bal, ok := h.mockBalances[skuID]
	if !ok {
		h.mu.Unlock()
		response.Error(c, http.StatusNotFound, sharedErrors.CodeSKUNotFound, "SKU库存未找到")
		return
	}
	if err := bal.Lock(qty); err != nil {
		h.mu.Unlock()
		response.Error(c, http.StatusOK, sharedErrors.CodeStockLockFailed, err.Error())
		return
	}
	h.mu.Unlock()

	response.Success(c, gin.H{
		"locked":          true,
		"available":       bal.Available(),
		"locked_quantity": bal.LockedQuantity,
	})
}

func (h *InventoryHandler) releaseInventory(c *gin.Context) {
	var req struct {
		LockKey  string `json:"lock_key" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	if h.mockEnabled {
		h.releaseInventoryMock(c, req.Quantity)
		return
	}

	ctx := c.Request.Context()
	if err := h.repo.ReleaseStock(ctx, req.LockKey, req.Quantity); err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"released": true})
}

func (h *InventoryHandler) releaseInventoryMock(c *gin.Context, qty int) {
	skuID := c.DefaultQuery("sku_id", "sku-001")

	h.mu.Lock()
	bal, ok := h.mockBalances[skuID]
	if !ok {
		h.mu.Unlock()
		response.Error(c, http.StatusNotFound, sharedErrors.CodeSKUNotFound, "SKU库存未找到")
		return
	}
	if err := bal.Release(qty); err != nil {
		h.mu.Unlock()
		response.Error(c, http.StatusOK, sharedErrors.CodeStockReleaseFailed, err.Error())
		return
	}
	h.mu.Unlock()

	response.Success(c, gin.H{"released": true, "available": bal.Available()})
}

func (h *InventoryHandler) deductInventory(c *gin.Context) {
	var req struct {
		LockKey string `json:"lock_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	if h.mockEnabled {
		h.deductInventoryMock(c)
		return
	}

	ctx := c.Request.Context()
	if err := h.repo.DeductStock(ctx, req.LockKey); err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"deducted": true})
}

func (h *InventoryHandler) deductInventoryMock(c *gin.Context) {
	skuID := c.DefaultQuery("sku_id", "sku-001")
	qtyStr := c.DefaultQuery("quantity", "10")

	qty, _ := strconv.Atoi(qtyStr)

	h.mu.Lock()
	bal, ok := h.mockBalances[skuID]
	if !ok {
		h.mu.Unlock()
		response.Error(c, http.StatusNotFound, sharedErrors.CodeSKUNotFound, "SKU库存未找到")
		return
	}
	if err := bal.Deduct(qty); err != nil {
		h.mu.Unlock()
		response.Error(c, http.StatusOK, sharedErrors.CodeStockDeductFailed, err.Error())
		return
	}
	h.mu.Unlock()

	response.Success(c, gin.H{"deducted": true, "total": bal.TotalQuantity})
}

func (h *InventoryHandler) listJournals(c *gin.Context) {
	if h.mockEnabled {
		h.listJournalsMock(c)
		return
	}

	tenantID := c.GetString("tenant_id")
	skuID := c.Query("sku_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	journals, total, err := h.repo.ListJournals(c.Request.Context(), tenantID, skuID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询流水失败: "+err.Error())
		return
	}
	response.PageSuccess(c, journals, total, page, pageSize)
}

func (h *InventoryHandler) listJournalsMock(c *gin.Context) {
	skuID := c.Query("sku_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	defaultSKU := "sku-001"
	if skuID == "" {
		defaultSKU = skuID
	}

	journals := []*domain.InventoryJournal{
		{ID: "j-1", SKUID: defaultSKU, ChangeType: "lock", ChangeQty: 10, BeforeTotal: 500, AfterTotal: 500, BeforeLocked: 0, AfterLocked: 10, BeforeAvail: 500, AfterAvail: 490},
	}
	response.PageSuccess(c, journals, int64(len(journals)), page, pageSize)
}
