package http

import (
	"net/http"
	"strconv"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/app"
	sharedEvents "github.com/Tangyd893/ERP-Go/backend/shared/events"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
	"github.com/gin-gonic/gin"
)

const fallbackMsg = "接口已联通，等待数据库迁移完成"

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	appService    *app.OrderAppService
	coordinator   *workflows.P4OutboundFlowCoordinator
	outboxStore   outbox.OutboxStore
	fallbackMode  bool
}

func NewOrderHandler(appService *app.OrderAppService) *OrderHandler {
	return &OrderHandler{
		appService:   appService,
		fallbackMode: appService == nil,
	}
}

func (h *OrderHandler) WithCoordinator(coordinator *workflows.P4OutboundFlowCoordinator) *OrderHandler {
	h.coordinator = coordinator
	return h
}

func (h *OrderHandler) WithOutboxStore(store outbox.OutboxStore) *OrderHandler {
	h.outboxStore = store
	return h
}

func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/orders", h.listOrders)
	router.GET("/orders/:id", h.getOrder)
	router.POST("/orders/audit", h.auditOrder)
	router.POST("/orders/abnormal", h.markAbnormal)
	router.POST("/orders/cancel", h.cancelOrder)
	router.POST("/fulfillment/outbound-shipped", h.outboundShipped)
	router.GET("/outbox/failed", h.listFailedOutbox)
	router.POST("/outbox/retry", h.retryOutbox)
}

func (h *OrderHandler) listOrders(c *gin.Context) {
	if h.fallbackMode {
		response.Success(c, []interface{}{})
		return
	}

	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	orders, total, err := h.appService.ListOrders(c.Request.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询订单失败: "+err.Error())
		return
	}
	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *OrderHandler) getOrder(c *gin.Context) {
	if h.fallbackMode {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeOrderNotFound, "订单不存在")
		return
	}

	order, err := h.appService.GetOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) auditOrder(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": fallbackMsg})
		return
	}

	var req struct {
		OrderID  string `json:"order_id" binding:"required"`
		Approved bool   `json:"approved"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	operator := c.GetString("username")
	if operator == "" {
		operator = "system"
	}

	if req.Approved {
		if err := h.appService.ApproveOrder(c.Request.Context(), req.OrderID, operator); err != nil {
			response.Error(c, http.StatusOK, sharedErrors.CodeOrderAuditFailed, err.Error())
			return
		}
		response.Success(c, gin.H{"approved": true})
	} else {
		if err := h.appService.CancelOrder(c.Request.Context(), req.OrderID, operator, "审核不通过"); err != nil {
			response.Error(c, http.StatusOK, sharedErrors.CodeOrderAuditFailed, err.Error())
			return
		}
		response.Success(c, gin.H{"approved": false, "cancelled": true})
	}
}

func (h *OrderHandler) markAbnormal(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": fallbackMsg})
		return
	}

	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Reason  string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	operator := c.GetString("username")
	if operator == "" {
		operator = "system"
	}

	if err := h.appService.MarkAbnormal(c.Request.Context(), req.OrderID, operator, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, gin.H{"abnormal": true})
}

func (h *OrderHandler) cancelOrder(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": fallbackMsg})
		return
	}

	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Reason  string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	operator := c.GetString("username")
	if operator == "" {
		operator = "system"
	}

	if err := h.appService.CancelOrder(c.Request.Context(), req.OrderID, operator, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, gin.H{"cancelled": true})
}

func (h *OrderHandler) outboundShipped(c *gin.Context) {
	if h.fallbackMode || h.coordinator == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，履约协调器未就绪"})
		return
	}

	var data workflows.OutboundShippedData
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	payload, err := outbox.NewEventPayload(sharedEvents.EventOutboundShipped, data)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		return
	}

	messageID := "ship-" + data.OutboundID
	if err := h.coordinator.HandleOutboundShipped(c.Request.Context(), messageID, payload); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "出库履约处理失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"processed": true, "order_id": data.OrderID})
}

func (h *OrderHandler) listFailedOutbox(c *gin.Context) {
	if h.outboxStore == nil {
		response.Error(c, http.StatusServiceUnavailable, sharedErrors.CodeInternalError, "Outbox 存储未就绪")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	messages, total, err := h.outboxStore.FetchFailed(c.Request.Context(), offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询失败消息失败: "+err.Error())
		return
	}

	type FailedItem struct {
		ID            int64  `json:"id"`
		AggregateID   string `json:"aggregate_id"`
		AggregateType string `json:"aggregate_type"`
		TenantID      string `json:"tenant_id"`
		EventType     string `json:"event_type"`
		RetryCount    int    `json:"retry_count"`
		CreatedAt     string `json:"created_at"`
		Status        string `json:"status"`
	}

	items := make([]FailedItem, 0, len(messages))
	for _, msg := range messages {
		items = append(items, FailedItem{
			ID:            msg.ID,
			AggregateID:   msg.AggregateID,
			AggregateType: msg.AggregateType,
			TenantID:      msg.TenantID,
			EventType:     msg.EventType,
			RetryCount:    msg.RetryCount,
			CreatedAt:     msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Status:        string(msg.Status),
		})
	}

	response.PageSuccess(c, items, total, page, pageSize)
}

func (h *OrderHandler) retryOutbox(c *gin.Context) {
	if h.outboxStore == nil {
		response.Error(c, http.StatusServiceUnavailable, sharedErrors.CodeInternalError, "Outbox 存储未就绪")
		return
	}

	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, sharedErrors.CodeInvalidParameter.Message())
		return
	}

	if err := h.outboxStore.Retry(c.Request.Context(), req.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "重试失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"retried": true, "id": req.ID})
}
