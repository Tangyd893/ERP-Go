package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	apperrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

// Server IAM HTTP 服务
type Server struct {
	authService *app.AuthService
	userService *app.UserService
	roleService *app.RoleService
	log         logger.Logger
	cfg         *config.Config
}

// NewServer 创建 IAM HTTP 服务
func NewServer(
	authService *app.AuthService,
	userService *app.UserService,
	roleService *app.RoleService,
	log logger.Logger,
	cfg *config.Config,
) *Server {
	return &Server{
		authService: authService,
		userService: userService,
		roleService: roleService,
		log:         log,
		cfg:         cfg,
	}
}

// RegisterRoutes 注册路由
func (s *Server) RegisterRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1/iam")
	{
		api.POST("/login", s.Login)
		api.POST("/refresh", s.RefreshToken)
		api.POST("/logout", s.authMiddleware(), s.Logout)

		api.GET("/user/info", s.authMiddleware(), s.GetUserInfo)
		api.POST("/check-permission", s.authMiddleware(), s.CheckPermission)

		userGroup := api.Group("/users")
		userGroup.Use(s.authMiddleware(), s.requirePermission("user:read"))
		{
			userGroup.GET("", s.ListUsers)
			userGroup.GET("/:id", s.GetUser)
			userGroup.POST("", s.requirePermission("user:create"), s.CreateUser)
			userGroup.PUT("/:id", s.requirePermission("user:update"), s.UpdateUser)
			userGroup.POST("/:id/roles", s.requirePermission("user:assign_role"), s.AssignUserRoles)
			userGroup.PUT("/:id/disable", s.requirePermission("user:disable"), s.DisableUser)
			userGroup.PUT("/:id/enable", s.requirePermission("user:enable"), s.EnableUser)
		}

		roleGroup := api.Group("/roles")
		roleGroup.Use(s.authMiddleware(), s.requirePermission("role:read"))
		{
			roleGroup.GET("", s.ListRoles)
			roleGroup.GET("/:id", s.GetRole)
			roleGroup.POST("", s.requirePermission("role:create"), s.CreateRole)
			roleGroup.PUT("/:id", s.requirePermission("role:update"), s.UpdateRole)
			roleGroup.POST("/:id/permissions", s.requirePermission("role:assign_perm"), s.AssignRolePermissions)
			roleGroup.DELETE("/:id", s.requirePermission("role:delete"), s.DeleteRole)
		}

		permGroup := api.Group("/permissions")
		permGroup.Use(s.authMiddleware())
		{
			permGroup.GET("", s.ListPermissions)
		}
	}
}

// authMiddleware 鉴权中间件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "未提供认证令牌")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, apperrors.CodeTokenInvalid, "令牌格式错误")
			return
		}

		tokenMgr := infra.NewJWTTokenManager(
			s.cfg.Server.Name, // placeholder
			2*time.Hour,
			7*24*time.Hour,
			"erp-go",
		)

		claims, err := tokenMgr.ValidateAccessToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, apperrors.CodeTokenInvalid, "令牌无效或已过期")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// requirePermission 权限校验中间件
func (s *Server) requirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			response.Error(c, http.StatusForbidden, apperrors.CodePermissionDenied, "权限不足")
			c.Abort()
			return
		}

		roleList := roles.([]string)
		hasPermission := false
		if contains(roleList, "super_admin") {
			hasPermission = true
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, apperrors.CodePermissionDenied, "权限不足: "+code)
			c.Abort()
			return
		}
		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Login 登录接口
func (s *Server) Login(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	result, err := s.authService.Login(c.Request.Context(), req.TenantID, req.Username, req.Password, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.BusinessError(c, err.(*apperrors.BusinessError))
		}
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新令牌接口
func (s *Server) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	result, err := s.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.BusinessError(c, err.(*apperrors.BusinessError))
		}
		return
	}

	response.Success(c, result)
}

// Logout 登出接口
func (s *Server) Logout(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")

	if err := s.authService.Logout(c.Request.Context(), tenantID, userID, username, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "登出失败")
		return
	}

	response.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
func (s *Server) GetUserInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	user, err := s.authService.GetUserInfo(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, apperrors.CodeNotFound, "用户不存在")
		return
	}

	response.Success(c, user)
}

// CheckPermission 权限校验接口
func (s *Server) CheckPermission(c *gin.Context) {
	var req struct {
		UserID         string `json:"user_id" binding:"required"`
		TenantID       string `json:"tenant_id" binding:"required"`
		PermissionCode string `json:"permission_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	allowed, err := s.authService.CheckPermission(c.Request.Context(), req.TenantID, req.UserID, req.PermissionCode)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "权限校验失败")
		return
	}

	response.Success(c, gin.H{"allowed": allowed})
}

// CreateUser 创建用户
func (s *Server) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	user, err := s.userService.CreateUser(c.Request.Context(), tenantID, req.Username, req.Password, req.Nickname, req.Email, req.Phone)
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "创建用户失败")
		}
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户
func (s *Server) UpdateUser(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	userID := c.Param("id")

	user, err := s.userService.UpdateUser(c.Request.Context(), tenantID, userID, req.Nickname, req.Email, req.Phone)
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "更新用户失败")
		}
		return
	}

	response.Success(c, user)
}

// AssignUserRoles 分配用户角色
func (s *Server) AssignUserRoles(c *gin.Context) {
	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	userID := c.Param("id")

	if err := s.userService.AssignRoles(c.Request.Context(), tenantID, userID, req.RoleIDs); err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "分配角色失败")
		}
		return
	}

	response.Success(c, nil)
}

// DisableUser 禁用用户
func (s *Server) DisableUser(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.Param("id")

	if err := s.userService.DisableUser(c.Request.Context(), tenantID, userID); err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "禁用用户失败")
		}
		return
	}

	response.Success(c, nil)
}

// EnableUser 启用用户
func (s *Server) EnableUser(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.Param("id")

	if err := s.userService.EnableUser(c.Request.Context(), tenantID, userID); err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "启用用户失败")
		}
		return
	}

	response.Success(c, nil)
}

// GetUser 获取用户详情
func (s *Server) GetUser(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.Param("id")

	user, err := s.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, apperrors.CodeNotFound, "用户不存在")
		return
	}

	response.Success(c, user)
}

// ListUsers 用户列表
func (s *Server) ListUsers(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	page := 1
	pageSize := 20

	users, total, err := s.userService.ListUsers(c.Request.Context(), tenantID, (page-1)*pageSize, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "查询用户列表失败")
		return
	}

	response.PageSuccess(c, users, total, page, pageSize)
}

// ListRoles 角色列表
func (s *Server) ListRoles(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	roles, total, err := s.roleService.ListRoles(c.Request.Context(), tenantID, 0, 100)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "查询角色列表失败")
		return
	}
	response.PageSuccess(c, roles, total, 1, 100)
}

// GetRole 获取角色详情
func (s *Server) GetRole(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	roleID := c.Param("id")

	role, err := s.roleService.GetRole(c.Request.Context(), tenantID, roleID)
	if err != nil {
		response.Error(c, http.StatusNotFound, apperrors.CodeNotFound, "角色不存在")
		return
	}

	response.Success(c, role)
}

// CreateRole 创建角色
func (s *Server) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	role, err := s.roleService.CreateRole(c.Request.Context(), tenantID, req.Name, req.Code, req.Description)
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "创建角色失败")
		}
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
func (s *Server) UpdateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	roleID := c.Param("id")

	role, err := s.roleService.UpdateRole(c.Request.Context(), tenantID, roleID, req.Name, req.Description)
	if err != nil {
		if bizErr, ok := err.(*apperrors.BusinessError); ok {
			response.BusinessError(c, bizErr)
		} else {
			response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "更新角色失败")
		}
		return
	}

	response.Success(c, role)
}

// AssignRolePermissions 给角色分配权限
func (s *Server) AssignRolePermissions(c *gin.Context) {
	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "参数无效")
		return
	}

	tenantID := c.GetString("tenant_id")
	roleID := c.Param("id")

	if err := s.roleService.AssignPermissions(c.Request.Context(), tenantID, roleID, req.PermissionIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "分配权限失败")
		return
	}

	response.Success(c, nil)
}

// DeleteRole 删除角色
func (s *Server) DeleteRole(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	roleID := c.Param("id")

	if err := s.roleService.DeleteRole(c.Request.Context(), tenantID, roleID); err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "删除角色失败")
		return
	}

	response.Success(c, nil)
}

// ListPermissions 权限列表
func (s *Server) ListPermissions(c *gin.Context) {
	perms, total, err := s.roleService.ListPermissions(c.Request.Context(), 0, 200)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, apperrors.CodeInternalError, "查询权限列表失败")
		return
	}
	response.PageSuccess(c, perms, total, 1, 200)
}
