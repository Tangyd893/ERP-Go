package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	router.POST("/arap/receivable", h.createReceivable)
	router.POST("/arap/payable", h.createPayable)
	router.GET("/costs", h.listCosts)
	router.POST("/costs", h.recordCost)
	router.GET("/profit", h.listProfit)
	router.POST("/profit/generate", h.generateProfit)
	router.GET("/exchange-rates", h.getExchangeRate)
	router.POST("/exchange-rates", h.setExchangeRate)
	router.GET("/journals", h.listJournals)
	router.POST("/journals", h.recordJournal)
}

// ── 结算 ────────────────────────────────────────────────

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
	var req struct {
		StoreID    string  `json:"store_id" binding:"required"`
		Platform   string  `json:"platform" binding:"required"`
		Period     string  `json:"period" binding:"required"`
		Currency   string  `json:"currency"`
		TotalSales float64 `json:"total_sales"`
		Refunds    float64 `json:"refunds"`
		Commission float64 `json:"commission"`
		FbaFee     float64 `json:"fba_fee"`
		OtherFee   float64 `json:"other_fee"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	bill, err := h.appService.ImportSettlement(c.Request.Context(), c.GetString("tenant_id"),
		req.StoreID, req.Platform, req.Period, req.Currency,
		req.TotalSales, req.Refunds, req.Commission, req.FbaFee, req.OtherFee)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "导入结算失败: "+err.Error())
		return
	}
	response.Success(c, bill)
}

// ── 应收应付 ────────────────────────────────────────────

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

func (h *FinanceHandler) createReceivable(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID  string  `json:"order_id" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	rec, err := h.appService.CreateReceivable(c.Request.Context(), c.GetString("tenant_id"), req.OrderID, req.Amount, req.Currency)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建应收失败: "+err.Error())
		return
	}
	response.Success(c, rec)
}

func (h *FinanceHandler) createPayable(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID  string  `json:"order_id" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	rec, err := h.appService.CreatePayable(c.Request.Context(), c.GetString("tenant_id"), req.OrderID, req.Amount, req.Currency)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建应付失败: "+err.Error())
		return
	}
	response.Success(c, rec)
}

// ── 成本 ────────────────────────────────────────────────

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

func (h *FinanceHandler) recordCost(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID  string  `json:"order_id" binding:"required"`
		SKUID    string  `json:"sku_id"`
		CostType string  `json:"cost_type" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	rec, err := h.appService.RecordCost(c.Request.Context(), c.GetString("tenant_id"), req.OrderID, req.SKUID, req.CostType, req.Amount, req.Currency)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "记录成本失败: "+err.Error())
		return
	}
	response.Success(c, rec)
}

// ── 利润 ────────────────────────────────────────────────

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

func (h *FinanceHandler) generateProfit(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID    string  `json:"order_id" binding:"required"`
		OrderNo    string  `json:"order_no" binding:"required"`
		SKUID      string  `json:"sku_id"`
		SKUCode    string  `json:"sku_code"`
		SaleAmount float64 `json:"sale_amount" binding:"required"`
		Currency   string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	report, err := h.appService.GenerateProfitReport(c.Request.Context(), c.GetString("tenant_id"),
		req.OrderID, req.OrderNo, req.SKUID, req.SKUCode, req.SaleAmount, req.Currency)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "生成利润报表失败: "+err.Error())
		return
	}
	response.Success(c, report)
}

// ── 汇率 ────────────────────────────────────────────────

func (h *FinanceHandler) getExchangeRate(c *gin.Context) {
	if h.fallbackMode { response.Success(c, gin.H{"rate": 1.0}); return }
	tenantID := c.GetString("tenant_id")
	from := c.DefaultQuery("from", "USD")
	to := c.DefaultQuery("to", "CNY")
	r, err := h.appService.GetExchangeRate(c.Request.Context(), tenantID, from, to)
	if err != nil {
		response.Success(c, gin.H{"rate": 1.0, "from": from, "to": to, "source": "default"})
		return
	}
	response.Success(c, r)
}

func (h *FinanceHandler) setExchangeRate(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		From   string  `json:"from" binding:"required"`
		To     string  `json:"to" binding:"required"`
		Rate   float64 `json:"rate" binding:"required"`
		Source string  `json:"source"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Source == "" { req.Source = "manual" }
	r, err := h.appService.SetExchangeRate(c.Request.Context(), c.GetString("tenant_id"), req.From, req.To, req.Rate, req.Source)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "设置汇率失败: "+err.Error())
		return
	}
	response.Success(c, r)
}

// ── 流水 ────────────────────────────────────────────────

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

func (h *FinanceHandler) recordJournal(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID        string  `json:"order_id" binding:"required"`
		ChangeType     string  `json:"change_type" binding:"required"`
		Amount         float64 `json:"amount" binding:"required"`
		BeforeAmount   float64 `json:"before_amount"`
		AfterAmount    float64 `json:"after_amount"`
		Currency       string  `json:"currency"`
		IdempotencyKey string  `json:"idempotency_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	if req.Currency == "" { req.Currency = "CNY" }
	if req.IdempotencyKey == "" { req.IdempotencyKey = fmt.Sprintf("JNL-%d", time.Now().UnixNano()) }
	j, err := h.appService.RecordJournal(c.Request.Context(), c.GetString("tenant_id"),
		req.OrderID, req.ChangeType, req.Amount, req.BeforeAmount, req.AfterAmount, req.Currency, req.IdempotencyKey)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "记录流水失败: "+err.Error())
		return
	}
	response.Success(c, j)
}
