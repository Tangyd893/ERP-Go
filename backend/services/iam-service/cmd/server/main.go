package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra/repository"
	httpiface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "erp-go-dev-secret-change-in-production"
	}

	var (
		authService *app.AuthService
		userService *app.UserService
		roleService *app.RoleService
	)

	srv := server.New(server.Options{
		Name:        "iam-service",
		DefaultPort: 8081,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			db, err := repository.NewDB(cfg)
			if err != nil {
				return nil, err
			}
			userRepo := repository.NewUserRepository(db)
			roleRepo := repository.NewRoleRepository(db)
			permRepo := repository.NewPermissionRepository(db)
			auditRepo := repository.NewAuditRepository(db)
			jwtMgr := infra.NewJWTTokenManager(jwtSecret, 2*time.Hour, 7*24*time.Hour, "erp-go")
			passHasher := infra.NewBcryptPasswordHasher(10)

			authService = app.NewAuthService(userRepo, roleRepo, jwtMgr, passHasher, auditRepo)
			userService = app.NewUserService(userRepo, roleRepo, passHasher, auditRepo)
			roleService = app.NewRoleService(roleRepo, permRepo, auditRepo)
			return db, nil
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			// IAM 需要 UserID 中间件（其他服务不需要）
			engine.Use(middleware.UserID())

			if authService != nil {
				cfg, _ := config.Load("")
				iamServer := httpiface.NewServer(authService, userService, roleService, log, cfg, jwtSecret)
				iamServer.RegisterRoutes(engine)
				log.Info("IAM 服务启动，数据库模式")
				return nil
			}
			api := engine.Group("/api/v1/iam")
			api.POST("/login", notImplYet)
			api.POST("/refresh", notImplYet)
			api.GET("/users", notImplYet)
			api.GET("/roles", notImplYet)
			api.GET("/permissions", notImplYet)
			log.Info("IAM 服务启动，占位模式")
			return nil
		},
	})
	srv.WaitShutdown()
}

func notImplYet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "接口已规划，数据库迁移完成后可用",
	})
}
