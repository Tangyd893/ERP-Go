package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	appService   *app.WarehouseAppService
	fallbackMode bool
}

func NewWarehouseHandler(appService *app.WarehouseAppService) *WarehouseHandler {
	return &WarehouseHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *WarehouseHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/outbounds", h.listOutbounds)
	router.POST("/outbounds", h.createOutbound)
	router.GET("/outbounds/:id", h.getOutbound)
	router.POST("/outbounds/:id/ship", h.confirmShip)
	router.GET("/pick-tasks", h.listPickTasks)
	router.POST("/pick/scan", h.pickScan)
	router.POST("/check/scan", h.checkScan)
	router.POST("/package", h.createPackage)
	router.POST("/weigh", h.weigh)
	router.GET("/locations", h.listLocations)
}

func (h *WarehouseHandler) listOutbounds(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	orders, total, err := h.appService.ListOutbounds(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询出库单失败: "+err.Error())
		return
	}
	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *WarehouseHandler) createOutbound(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID     string `json:"order_id" binding:"required"`
		OrderNo     string `json:"order_no" binding:"required"`
		WarehouseID string `json:"warehouse_id" binding:"required"`
		Items       []struct {
			SKUID    string `json:"sku_id"`
			SKUCode  string `json:"sku_code"`
			SKUName  string `json:"sku_name"`
			Quantity int    `json:"quantity"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	now := time.Now()
	items := make([]*domain.OutboundItem, 0, len(req.Items))
	for i, it := range req.Items {
		items = append(items, &domain.OutboundItem{
			ID: fmt.Sprintf("OI%d-%d", now.UnixNano(), i),
			SKUID: it.SKUID, SKUCode: it.SKUCode, SKUName: it.SKUName, Quantity: it.Quantity,
		})
	}
	order := &domain.OutboundOrder{
		ID: fmt.Sprintf("OB%d", now.UnixNano()), TenantID: c.GetString("tenant_id"),
		OrderID: req.OrderID, OrderNo: req.OrderNo, WarehouseID: req.WarehouseID,
		Status: domain.OutboundPicking, Items: items, CreatedAt: now, UpdatedAt: now,
	}
	if err := h.appService.CreateOutbound(c.Request.Context(), order); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建出库单失败: "+err.Error())
		return
	}
	response.Success(c, order)
}

func (h *WarehouseHandler) getOutbound(c *gin.Context) {
	if h.fallbackMode { response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "出库单不存在"); return }
	order, err := h.appService.GetOutbound(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "出库单不存在")
		return
	}
	response.Success(c, order)
}

func (h *WarehouseHandler) confirmShip(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		TrackingNo string `json:"tracking_no"`
		Carrier    string `json:"carrier"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.appService.ConfirmShip(c.Request.Context(), c.Param("id"), req.TrackingNo, req.Carrier); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "出库确认失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"shipped": true})
}

func (h *WarehouseHandler) listPickTasks(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	outboundID := c.Query("outbound_id")
	if outboundID == "" {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "缺少outbound_id")
		return
	}
	tasks, err := h.appService.ListPickTasks(c.Request.Context(), outboundID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询拣货任务失败: "+err.Error())
		return
	}
	response.Success(c, tasks)
}

func (h *WarehouseHandler) pickScan(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		TaskID string `json:"task_id" binding:"required"`
		Qty    int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	if err := h.appService.PickScan(c.Request.Context(), req.TaskID, req.Qty); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "拣货扫码失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"picked": true})
}

func (h *WarehouseHandler) checkScan(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OutboundID string `json:"outbound_id" binding:"required"`
		SKUID      string `json:"sku_id"`
		Qty        int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	if err := h.appService.UpdateOutboundStatus(c.Request.Context(), req.OutboundID, string(domain.OutboundChecked)); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "复核失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"checked": true})
}

func (h *WarehouseHandler) createPackage(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OutboundID string  `json:"outbound_id" binding:"required"`
		Weight     float64 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	if err := h.appService.UpdateOutboundStatus(c.Request.Context(), req.OutboundID, string(domain.OutboundPacked)); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "打包失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"packed": true})
}

func (h *WarehouseHandler) weigh(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OutboundID string  `json:"outbound_id" binding:"required"`
		Weight     float64 `json:"weight" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	if err := h.appService.UpdateOutboundStatus(c.Request.Context(), req.OutboundID, string(domain.OutboundWeighed)); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "称重失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"weighed": true, "weight": req.Weight})
}

func (h *WarehouseHandler) listLocations(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	whs, err := h.appService.ListWarehouses(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询仓库失败: "+err.Error())
		return
	}
	response.Success(c, whs)
}
