package http

import (
	"net/http"

	"github.com/Tangyd893/ERP-Go/backend/services/report-service/internal/app"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	appService *app.ReportAppService
}

func NewReportHandler(appService *app.ReportAppService) *ReportHandler {
	return &ReportHandler{appService: appService}
}

func (h *ReportHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", h.dashboard)
	router.GET("/sales", h.salesReport)
	router.GET("/inventory-turnover", h.inventoryTurnover)
	router.GET("/warehouse-efficiency", h.warehouseEfficiency)
	router.GET("/profit-summary", h.profitSummary)
}

func (h *ReportHandler) dashboard(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	data, err := h.appService.GetDashboard(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询看板数据失败: "+err.Error())
		return
	}
	response.Success(c, data)
}

func (h *ReportHandler) salesReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	period := c.DefaultQuery("period", "monthly")
	report, err := h.appService.GetSalesReport(c.Request.Context(), tenantID, period)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询销售报表失败: "+err.Error())
		return
	}
	response.Success(c, report)
}

func (h *ReportHandler) inventoryTurnover(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	data, err := h.appService.GetInventoryTurnover(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询库存周转失败: "+err.Error())
		return
	}
	response.Success(c, data)
}

func (h *ReportHandler) warehouseEfficiency(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	data, err := h.appService.GetWarehouseEfficiency(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询仓储效率失败: "+err.Error())
		return
	}
	response.Success(c, data)
}

func (h *ReportHandler) profitSummary(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	period := c.DefaultQuery("period", "monthly")
	report, err := h.appService.GetProfitSummary(c.Request.Context(), tenantID, period)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询利润汇总失败: "+err.Error())
		return
	}
	response.Success(c, report)
}
