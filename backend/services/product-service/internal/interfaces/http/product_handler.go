package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

// ProductHandler 商品 HTTP 处理器
type ProductHandler struct {
	appService   *app.ProductAppService
	fallbackMode bool
}

func NewProductHandler(appService *app.ProductAppService) *ProductHandler {
	return &ProductHandler{
		appService:   appService,
		fallbackMode: appService == nil,
	}
}

func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/spus", h.listSPUs)
	router.GET("/skus", h.listSKUs)
	router.POST("/skus", h.createSKU)
	router.GET("/skus/:code", h.getSKU)
	router.GET("/sku-mappings", h.listSKUMappings)
	router.POST("/sku-mappings", h.createSKUMapping)
}

func (h *ProductHandler) listSPUs(c *gin.Context) {
	if h.fallbackMode {
		response.PageSuccess(c, []interface{}{}, 0, 1, 20)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	response.PageSuccess(c, []interface{}{}, 0, page, pageSize)
}

func (h *ProductHandler) listSKUs(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	skus, total, err := h.appService.ListSKUs(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询SKU失败: "+err.Error())
		return
	}
	response.PageSuccess(c, skus, total, page, pageSize)
}

func (h *ProductHandler) createSKU(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		SPUID    string  `json:"spu_id"`
		Code     string  `json:"code" binding:"required"`
		Barcode  string  `json:"barcode"`
		Weight   float64 `json:"weight"`
		Length   float64 `json:"length"`
		Width    float64 `json:"width"`
		Height   float64 `json:"height"`
		SalePrice float64 `json:"sale_price"`
		Currency string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	sku := &domain.SKU{
		ID:       fmt.Sprintf("SKU%d", time.Now().UnixNano()),
		TenantID: c.GetString("tenant_id"),
		SPUID:    req.SPUID,
		Code:     req.Code,
		Barcode:  req.Barcode,
		Weight:   req.Weight,
		Length:   req.Length,
		Width:    req.Width,
		Height:   req.Height,
		SalePrice: req.SalePrice,
		Currency: req.Currency,
		Status:   "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.appService.CreateSKU(c.Request.Context(), sku); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建SKU失败: "+err.Error())
		return
	}
	response.Success(c, sku)
}

func (h *ProductHandler) getSKU(c *gin.Context) {
	if h.fallbackMode {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeSKUNotFound, "SKU不存在")
		return
	}

	tenantID := c.GetString("tenant_id")
	sku, err := h.appService.GetSKU(c.Request.Context(), tenantID, c.Param("code"))
	if err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, sku)
}

func (h *ProductHandler) listSKUMappings(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	tenantID := c.GetString("tenant_id")
	storeID := c.Query("store_id")

	mappings := make([]*domain.PlatformSKU, 0)
	if storeID != "" {
		mapping, err := h.appService.GetPlatformSKU(c.Request.Context(), tenantID, storeID, c.Query("platform_code"))
		if err == nil {
			mappings = append(mappings, mapping)
		}
	}
	response.Success(c, mappings)
}

func (h *ProductHandler) createSKUMapping(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		SKUID        string `json:"sku_id" binding:"required"`
		StoreID      string `json:"store_id" binding:"required"`
		PlatformCode string `json:"platform_code" binding:"required"`
		PlatformSKU  string `json:"platform_sku" binding:"required"`
		ASIN         string `json:"asin"`
		FNSKU        string `json:"fnsku"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	mapping := &domain.PlatformSKU{
		ID:           fmt.Sprintf("PSK%d", time.Now().UnixNano()),
		TenantID:     c.GetString("tenant_id"),
		SKUID:        req.SKUID,
		StoreID:      req.StoreID,
		PlatformCode: req.PlatformCode,
		PlatformSKU:  req.PlatformSKU,
		ASIN:         req.ASIN,
		FNSKU:        req.FNSKU,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := h.appService.MapPlatformSKU(c.Request.Context(), mapping); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建SKU映射失败: "+err.Error())
		return
	}
	response.Success(c, mapping)
}
