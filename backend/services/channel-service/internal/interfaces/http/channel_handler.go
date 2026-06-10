package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

// ChannelHandler 渠道 HTTP 处理器
type ChannelHandler struct {
	appService   *app.ChannelAppService
	fallbackMode bool
}

func NewChannelHandler(appService *app.ChannelAppService) *ChannelHandler {
	return &ChannelHandler{
		appService:   appService,
		fallbackMode: appService == nil,
	}
}

func (h *ChannelHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/stores", h.listStores)
	router.POST("/stores", h.createStore)
	router.POST("/orders/import", h.importOrders)
	router.GET("/import-tasks", h.listImportTasks)
	router.GET("/sync-tasks", h.listSyncTasks)
}

func (h *ChannelHandler) listStores(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	tenantID := c.GetString("tenant_id")
	stores, err := h.appService.ListStores(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询店铺失败: "+err.Error())
		return
	}
	response.Success(c, stores)
}

func (h *ChannelHandler) createStore(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		Name         string `json:"name" binding:"required"`
		PlatformCode string `json:"platform_code" binding:"required"`
		Site         string `json:"site"`
		StoreCode    string `json:"store_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	store := &domain.Store{
		ID:           fmt.Sprintf("ST%d", time.Now().UnixNano()),
		TenantID:     c.GetString("tenant_id"),
		PlatformCode: req.PlatformCode,
		Site:         req.Site,
		Name:         req.Name,
		StoreCode:    req.StoreCode,
		Status:       domain.StoreStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := h.appService.CreateStore(c.Request.Context(), store); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建店铺失败: "+err.Error())
		return
	}
	response.Success(c, store)
}

func (h *ChannelHandler) importOrders(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		StoreID  string `json:"store_id" binding:"required"`
		FileName string `json:"file_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	task := &domain.OrderImportTask{
		ID:             fmt.Sprintf("IM%d", time.Now().UnixNano()),
		TenantID:       c.GetString("tenant_id"),
		StoreID:        req.StoreID,
		ImportType:     "csv",
		FileName:       req.FileName,
		IdempotencyKey: fmt.Sprintf("%s-%d", req.StoreID, time.Now().Unix()),
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := h.appService.CreateImportTask(c.Request.Context(), task); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建导入任务失败: "+err.Error())
		return
	}
	response.Success(c, task)
}

func (h *ChannelHandler) listImportTasks(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	// 占位：导入任务列表通过幂等键查询，后续扩展为完整列表
	key := c.Query("idempotency_key")
	if key != "" {
		task, err := h.appService.GetImportTask(c.Request.Context(), key)
		if err != nil {
			response.Error(c, http.StatusNotFound, sharedErrors.CodeNotFound, "导入任务不存在")
			return
		}
		response.Success(c, []*domain.OrderImportTask{task})
		return
	}
	response.PageSuccess(c, []interface{}{}, 0, 1, 20)
}

func (h *ChannelHandler) listSyncTasks(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	response.PageSuccess(c, []interface{}{}, 0, page, pageSize)
}
