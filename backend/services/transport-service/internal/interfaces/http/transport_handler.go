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
	router.POST("/match-carrier", h.matchCarrier)
	router.POST("/shipments", h.createShipment)
	router.POST("/labels", h.createLabel)
	router.POST("/labels/generate", h.generateLabel)
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
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
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
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
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

// matchCarrier 物流渠道匹配：按重量+目的地匹配最优物流产品
func (h *TransportHandler) matchCarrier(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		Weight  float64 `json:"weight" binding:"required"`
		Country string  `json:"country" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	tenantID := c.GetString("tenant_id")
	result, err := h.appService.MatchCarrier(c.Request.Context(), tenantID, req.Weight, req.Country)
	if err != nil {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeInvalidParameter, "物流匹配失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"matched":         true,
		"rule_id":         result.RuleID,
		"rule_name":       result.RuleName,
		"carrier_service": result.CarrierService,
	})
}

// generateLabel 调用适配器生成面单
func (h *TransportHandler) generateLabel(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	var req struct {
		ShipmentID string `json:"shipment_id" binding:"required"`
		OrderNo    string `json:"order_no"`
		Weight     float64 `json:"weight"`
		Country    string `json:"country"`
		Address    struct {
			Name       string `json:"name"`
			Phone      string `json:"phone"`
			Country    string `json:"country"`
			State      string `json:"state"`
			City       string `json:"city"`
			District   string `json:"district"`
			StreetLine string `json:"street_line"`
			PostalCode string `json:"postal_code"`
		} `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}
	addr := app.AddressInfo{
		Name: req.Address.Name, Phone: req.Address.Phone, Country: req.Address.Country,
		State: req.Address.State, City: req.Address.City, District: req.Address.District,
		StreetLine: req.Address.StreetLine, PostalCode: req.Address.PostalCode,
	}
	resp, err := h.appService.GenerateLabel(c.Request.Context(), req.ShipmentID, req.OrderNo, req.Country, req.Weight, addr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "生成面单失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"created":      true,
		"tracking_no":  resp.TrackingNo,
		"label_url":    resp.LabelURL,
		"carrier_code": resp.CarrierCode,
		"service_code": resp.ServiceCode,
	})
}
