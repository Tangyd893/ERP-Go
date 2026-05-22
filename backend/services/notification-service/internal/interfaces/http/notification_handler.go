package http

import (
	"net/http"
	"strconv"

	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/app"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	appService   *app.NotificationAppService
	fallbackMode bool
}

func NewNotificationHandler(appService *app.NotificationAppService) *NotificationHandler {
	return &NotificationHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/list", h.listNotifications)
	router.PUT("/read", h.markRead)
	router.GET("/unread-count", h.unreadCount)
}

func (h *NotificationHandler) listNotifications(c *gin.Context) {
	if h.fallbackMode { response.Success(c, []interface{}{}); return }
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	list, total, err := h.appService.ListNotifications(c.Request.Context(), tenantID, userID, offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询通知失败: "+err.Error())
		return
	}
	response.PageSuccess(c, list, total, page, pageSize)
}

func (h *NotificationHandler) markRead(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通"}); return }
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	if err := h.appService.MarkAllRead(c.Request.Context(), tenantID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "标记已读失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"read": true})
}

func (h *NotificationHandler) unreadCount(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "data": 0}); return }
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	count, err := h.appService.GetUnreadCount(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询未读数失败: "+err.Error())
		return
	}
	response.Success(c, count)
}
