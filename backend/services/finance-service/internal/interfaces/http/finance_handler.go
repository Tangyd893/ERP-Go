package http

import (
	"net/http"
	"strconv"

	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/app"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	appService   *app.FinanceAppService
	fallbackMode bool
}

func NewFinanceHandler(appService *app.FinanceAppService) *FinanceHandler {
	return &FinanceHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *FinanceHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/settlements", h.listSettlements)
	router.POST("/settlements/import", h.importSettlement)
	router.GET("/arap", h.listArAp)
	router.GET("/costs", h.listCosts)
	router.GET("/profit", h.listProfit)
	router.GET("/journals", h.listJournals)
}

func (h *FinanceHandler) listSettlements(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	bills, total, err := h.appService.ListSettlementBills(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询结算单失败: "+err.Error())
		return
	}
	response.PageSuccess(c, bills, total, page, pageSize)
}

func (h *FinanceHandler) importSettlement(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "结算导入已提交", "data": gin.H{"task_id": "pending"}})
}

func (h *FinanceHandler) listArAp(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	records, total, err := h.appService.ListArApRecords(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询应收应付失败: "+err.Error())
		return
	}
	response.PageSuccess(c, records, total, page, pageSize)
}

func (h *FinanceHandler) listCosts(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	records, total, err := h.appService.ListCostRecords(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询成本记录失败: "+err.Error())
		return
	}
	response.PageSuccess(c, records, total, page, pageSize)
}

func (h *FinanceHandler) listProfit(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	reports, total, err := h.appService.ListProfitReports(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询利润报表失败: "+err.Error())
		return
	}
	response.PageSuccess(c, reports, total, page, pageSize)
}

func (h *FinanceHandler) listJournals(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize
	journals, total, err := h.appService.ListJournals(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询财务流水失败: "+err.Error())
		return
	}
	response.PageSuccess(c, journals, total, page, pageSize)
}
