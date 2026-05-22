package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

// TenantHandler 租户 HTTP 处理器
type TenantHandler struct {
	appService  *app.TenantAppService
	fallbackMode bool
}

// NewTenantHandler 创建租户处理器，appService 为 nil 时使用占位模式
func NewTenantHandler(appService *app.TenantAppService) *TenantHandler {
	return &TenantHandler{
		appService:   appService,
		fallbackMode: appService == nil,
	}
}

// RegisterRoutes 注册路由
func (h *TenantHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/tenants", h.listTenants)
	router.GET("/tenants/:id", h.getTenant)
	router.POST("/tenants", h.createTenant)
	router.PUT("/tenants/:id", h.updateTenant)

	router.GET("/organizations", h.listOrganizations)
	router.POST("/organizations", h.createOrganization)

	router.GET("/departments", h.listDepartments)
	router.POST("/departments", h.createDepartment)

	router.GET("/positions", h.listPositions)
	router.POST("/positions", h.createPosition)
}

func (h *TenantHandler) listTenants(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}, "total": 0})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	tenants, total, err := h.appService.ListTenants(c.Request.Context(), offset, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询租户失败: "+err.Error())
		return
	}
	response.PageSuccess(c, tenants, total, page, pageSize)
}

func (h *TenantHandler) getTenant(c *gin.Context) {
	if h.fallbackMode {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeTenantNotFound, "租户不存在")
		return
	}

	tenant, err := h.appService.GetTenant(c.Request.Context(), c.Param("id"))
	if err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, tenant)
}

func (h *TenantHandler) createTenant(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		Name         string `json:"name" binding:"required"`
		Code         string `json:"code" binding:"required"`
		ContactName  string `json:"contact_name"`
		ContactEmail string `json:"contact_email"`
		ContactPhone string `json:"contact_phone"`
		QuotaUsers   int    `json:"quota_users"`
		QuotaOrders  int    `json:"quota_orders"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenant := &domain.Tenant{
		ID:           fmt.Sprintf("TN%d", time.Now().UnixNano()),
		Name:         req.Name,
		Code:         req.Code,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       domain.TenantStatusActive,
		QuotaUsers:   req.QuotaUsers,
		QuotaOrders:  req.QuotaOrders,
	}
	if err := h.appService.CreateTenant(c.Request.Context(), tenant); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建租户失败: "+err.Error())
		return
	}
	response.Success(c, tenant)
}

func (h *TenantHandler) updateTenant(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		Name         string `json:"name"`
		ContactName  string `json:"contact_name"`
		ContactEmail string `json:"contact_email"`
		ContactPhone string `json:"contact_phone"`
		Status       string `json:"status"`
		QuotaUsers   int    `json:"quota_users"`
		QuotaOrders  int    `json:"quota_orders"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenant, err := h.appService.GetTenant(c.Request.Context(), c.Param("id"))
	if err != nil {
		if bizErr, ok := err.(*sharedErrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, err.Error())
		}
		return
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.ContactName != "" {
		tenant.ContactName = req.ContactName
	}
	if req.ContactEmail != "" {
		tenant.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		tenant.ContactPhone = req.ContactPhone
	}
	if req.Status != "" {
		tenant.Status = domain.TenantStatus(req.Status)
	}
	if req.QuotaUsers > 0 {
		tenant.QuotaUsers = req.QuotaUsers
	}
	if req.QuotaOrders > 0 {
		tenant.QuotaOrders = req.QuotaOrders
	}

	if err := h.appService.UpdateTenant(c.Request.Context(), tenant); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "更新租户失败: "+err.Error())
		return
	}
	response.Success(c, tenant)
}

func (h *TenantHandler) listOrganizations(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	tenantID := c.GetString("tenant_id")
	orgs, err := h.appService.ListOrganizations(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询组织失败: "+err.Error())
		return
	}
	response.Success(c, orgs)
}

func (h *TenantHandler) createOrganization(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		ParentID  string `json:"parent_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	org := &domain.Organization{
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    "active",
	}
	if err := h.appService.CreateOrganization(c.Request.Context(), org); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建组织失败: "+err.Error())
		return
	}
	response.Success(c, org)
}

func (h *TenantHandler) listDepartments(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	orgID := c.Query("org_id")
	if orgID == "" {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "缺少组织ID参数")
		return
	}

	depts, err := h.appService.ListDepartments(c.Request.Context(), orgID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询部门失败: "+err.Error())
		return
	}
	response.Success(c, depts)
}

func (h *TenantHandler) createDepartment(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		OrgID     string `json:"org_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		ParentID  string `json:"parent_id"`
		ManagerID string `json:"manager_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	dept := &domain.Department{
		OrgID:     req.OrgID,
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		ManagerID: req.ManagerID,
		SortOrder: req.SortOrder,
		Status:    "active",
	}
	if err := h.appService.CreateDepartment(c.Request.Context(), dept); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建部门失败: "+err.Error())
		return
	}
	response.Success(c, dept)
}

func (h *TenantHandler) listPositions(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	deptID := c.Query("dept_id")
	if deptID == "" {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "缺少部门ID参数")
		return
	}

	positions, err := h.appService.ListPositions(c.Request.Context(), deptID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "查询岗位失败: "+err.Error())
		return
	}
	response.Success(c, positions)
}

func (h *TenantHandler) createPosition(c *gin.Context) {
	if h.fallbackMode {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，等待数据库迁移完成"})
		return
	}

	var req struct {
		DeptID    string `json:"dept_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "参数无效")
		return
	}

	pos := &domain.Position{
		DeptID:    req.DeptID,
		Name:      req.Name,
		Code:      req.Code,
		SortOrder: req.SortOrder,
		Status:    "active",
	}
	if err := h.appService.CreatePosition(c.Request.Context(), pos); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "创建岗位失败: "+err.Error())
		return
	}
	response.Success(c, pos)
}
