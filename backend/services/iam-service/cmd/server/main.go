package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra/repository"
	grpciface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/grpc"
	httpiface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	redisx "github.com/Tangyd893/ERP-Go/backend/shared/redis"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = config.DefaultJWTSecret
	}

	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	if config.RequiresSecureConfig() {
		if err := config.ValidateProduction(jwtSecret, cfg.Database.Password); err != nil {
			panic(fmt.Sprintf("安全配置校验失败: %v", err))
		}
	}

	var (
		authService *app.AuthService
		userService *app.UserService
		roleService *app.RoleService
		redisClient *redisx.Client
	)

	stopGRPC := func() {}
	stopGRPCPtr := &stopGRPC

	srv := server.New(server.Options{
		Name:        "iam-service",
		DefaultPort: 8081,
		OnShutdown: func() {
			(*stopGRPCPtr)()
			if redisClient != nil {
				_ = redisClient.Close()
			}
		},
		InitDB: func(dbCfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			db, err := repository.NewDB(dbCfg)
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

			if rc, rerr := redisx.New(cfg.Redis); rerr != nil {
				log.Warnf("Redis 不可用，跳过会话缓存: %v", rerr)
			} else if rc != nil {
				redisClient = rc
				authService.WithSessionStore(redisx.NewSessionStore(rc))
				if store := redisx.NewRefreshRevokeStore(rc); store != nil {
					jwtMgr.WithRevocationStore(store)
					log.Info("Redis refresh 吊销已启用")
				}
				log.Info("Redis 会话缓存已启用")
			}
			return db, nil
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			engine.Use(middleware.UserID())

			if authService != nil {
				loaded, _ := config.Load("")
				iamServer := httpiface.NewServer(authService, userService, roleService, log, loaded, jwtSecret)
				iamServer.RegisterRoutes(engine)
				*stopGRPCPtr = grpciface.TryStart(authService, log)
				log.Info("IAM 服务启动，数据库模式")
				return nil
			}
			if config.IsDevelopment() {
				api := engine.Group("/api/v1/iam")
				api.POST("/login", notImplYet)
				api.POST("/refresh", notImplYet)
				api.GET("/users", notImplYet)
				api.GET("/roles", notImplYet)
				api.GET("/permissions", notImplYet)
				log.Info("IAM 服务启动，占位模式")
			}
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
