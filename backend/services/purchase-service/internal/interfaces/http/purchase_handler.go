package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	appService   *app.PurchaseAppService
	fallbackMode bool
}

func NewPurchaseHandler(appService *app.PurchaseAppService) *PurchaseHandler {
	return &PurchaseHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *PurchaseHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/suppliers", h.listSuppliers)
	router.POST("/suppliers", h.createSupplier)
	router.GET("/orders", h.listOrders)
	router.POST("/orders", h.createOrder)
	router.GET("/orders/:id", h.getOrder)
	router.POST("/orders/:id/submit", h.submitOrder)
	router.POST("/orders/:id/approve", h.approveOrder)
	router.POST("/orders/:id/ordered", h.markOrdered)
	router.POST("/orders/:id/cancel", h.cancelOrder)
	router.POST("/orders/:id/receive", h.receiveItem)
	router.GET("/inbound", h.listInbound)
	router.GET("/inbound/:id", h.getInbound)
	router.POST("/inbound/:id/qa", h.startQA)
	router.POST("/inbound/:id/qa-items", h.qaItem)
	router.POST("/inbound/:id/complete", h.completeInbound)
	router.POST("/inbound/:id/return", h.returnRejected)
}

func (h *PurchaseHandler) listSuppliers(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	suppliers, err := h.appService.ListSuppliers(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询供应商失败: "+err.Error())
		return
	}
	response.Success(c, suppliers)
}

func (h *PurchaseHandler) createSupplier(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		Name         string `json:"name" binding:"required"`
		Code         string `json:"code" binding:"required"`
		ContactName  string `json:"contact_name"`
		ContactPhone string `json:"contact_phone"`
		Email        string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	s := &domain.Supplier{
		ID: fmt.Sprintf("SP%d", time.Now().UnixNano()), TenantID: c.GetString("tenant_id"),
		Name: req.Name, Code: req.Code, ContactName: req.ContactName, ContactPhone: req.ContactPhone,
		Email: req.Email, Status: "active", CreatedAt: time.Now(),
	}
	if err := h.appService.CreateSupplier(c.Request.Context(), s); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建供应商失败: "+err.Error())
		return
	}
	response.Success(c, s)
}

func (h *PurchaseHandler) listOrders(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	orders, total, err := h.appService.ListPurchaseOrders(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询采购单失败: "+err.Error())
		return
	}
	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *PurchaseHandler) createOrder(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		SupplierID string  `json:"supplier_id" binding:"required"`
		TotalAmount float64 `json:"total_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	order := &domain.PurchaseOrder{
		ID: fmt.Sprintf("PO%d", time.Now().UnixNano()), TenantID: c.GetString("tenant_id"),
		SupplierID: req.SupplierID, OrderNo: fmt.Sprintf("PO-%d", time.Now().Unix()),
		Status: domain.PurchaseDraft, TotalAmount: req.TotalAmount, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := h.appService.CreatePurchaseOrder(c.Request.Context(), order); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建采购单失败: "+err.Error())
		return
	}
	response.Success(c, order)
}

func (h *PurchaseHandler) listInbound(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	orders, total, err := h.appService.ListInboundOrders(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询入库单失败: "+err.Error())
		return
	}
	response.PageSuccess(c, orders, total, page, pageSize)
}

// ── 采购单工作流 ────────────────────────────────────────

func (h *PurchaseHandler) getOrder(c *gin.Context) {
	if h.fallbackMode { response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "采购单不存在"); return }
	order, err := h.appService.GetPurchaseOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "采购单不存在")
		return
	}
	response.Success(c, order)
}

func (h *PurchaseHandler) submitOrder(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.SubmitOrder(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "提交失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"submitted": true})
}

func (h *PurchaseHandler) approveOrder(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.ApproveOrder(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "审核失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"approved": true})
}

func (h *PurchaseHandler) markOrdered(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.MarkOrdered(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "下单失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"ordered": true})
}

func (h *PurchaseHandler) cancelOrder(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.CancelOrder(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "取消失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"cancelled": true})
}

func (h *PurchaseHandler) receiveItem(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		ItemID      string `json:"item_id" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
		WarehouseID string `json:"warehouse_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	inbound, err := h.appService.ReceiveItem(c.Request.Context(), c.Param("id"), req.ItemID, req.WarehouseID, req.Quantity)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "收货失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"received": true, "inbound": inbound})
}

// ── 入库单工作流 ────────────────────────────────────────

func (h *PurchaseHandler) getInbound(c *gin.Context) {
	if h.fallbackMode { response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "入库单不存在"); return }
	in, err := h.appService.GetInboundOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "入库单不存在")
		return
	}
	response.Success(c, in)
}

func (h *PurchaseHandler) startQA(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.StartQA(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "开始质检失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"qa_started": true})
}

func (h *PurchaseHandler) qaItem(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		ItemID   string `json:"item_id" binding:"required"`
		Passed   int    `json:"passed"`
		Rejected int    `json:"rejected"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if err := h.appService.QAItem(c.Request.Context(), c.Param("id"), req.ItemID, req.Passed, req.Rejected); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "质检失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"qa_done": true})
}

func (h *PurchaseHandler) completeInbound(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.CompleteInbound(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "完成入库失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"inbound_completed": true})
}

func (h *PurchaseHandler) returnRejected(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	if err := h.appService.ReturnRejectedItems(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "退货失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"returned": true})
}
