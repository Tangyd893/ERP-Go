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
	router.GET("/inbound", h.listInbound)
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
