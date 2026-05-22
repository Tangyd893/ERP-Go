package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type TransportHandler struct {
	appService   *app.TransportAppService
	fallbackMode bool
}

func NewTransportHandler(appService *app.TransportAppService) *TransportHandler {
	return &TransportHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *TransportHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/carriers", h.listCarriers)
	router.GET("/rules", h.listRules)
	router.POST("/shipments", h.createShipment)
	router.POST("/labels", h.createLabel)
	router.GET("/tracking", h.getTracking)
}

func (h *TransportHandler) listCarriers(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	carriers, err := h.appService.ListCarriers(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询物流商失败: "+err.Error())
		return
	}
	response.Success(c, carriers)
}

func (h *TransportHandler) listRules(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	response.PageSuccess(c, []interface{}{}, 0, page, pageSize)
}

func (h *TransportHandler) createShipment(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		OrderID    string  `json:"order_id" binding:"required"`
		OutboundID string  `json:"outbound_id" binding:"required"`
		CarrierCode string `json:"carrier_code" binding:"required"`
		Weight     float64 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	s := &domain.Shipment{
		ID: fmt.Sprintf("SH%d", time.Now().UnixNano()), TenantID: c.GetString("tenant_id"),
		OrderID: req.OrderID, OutboundID: req.OutboundID, CarrierCode: req.CarrierCode,
		Weight: req.Weight, Status: domain.ShipmentPending, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := h.appService.CreateShipment(c.Request.Context(), s); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建发运单失败: "+err.Error())
		return
	}
	response.Success(c, s)
}

func (h *TransportHandler) createLabel(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		ShipmentID string `json:"shipment_id" binding:"required"`
		TrackingNo string `json:"tracking_no" binding:"required"`
		LabelURL   string `json:"label_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}
	if err := h.appService.CreateLabel(c.Request.Context(), req.ShipmentID, req.TrackingNo, req.LabelURL); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建面单失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"labeled": true, "tracking_no": req.TrackingNo})
}

func (h *TransportHandler) getTracking(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	shipmentID := c.Query("shipment_id")
	if shipmentID == "" {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "缺少shipment_id")
		return
	}
	s, err := h.appService.GetShipment(c.Request.Context(), shipmentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询轨迹失败: "+err.Error())
		return
	}
	response.Success(c, s)
}
