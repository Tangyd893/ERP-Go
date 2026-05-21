package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra/repository"
	httpiface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	cfg.Server.Name = "iam-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8081
	}

	log := logger.New(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Server.Name,
		os.Getenv("ENVIRONMENT"),
	)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "erp-go-dev-secret-change-in-production"
	}

	var (
		db          *gorm.DB
		authService *app.AuthService
		userService *app.UserService
		roleService *app.RoleService
	)

	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用占位模式: %v", dbErr)
	} else {
		log.Info("数据库连接成功")
		db = database

		userRepo := repository.NewUserRepository(db)
		roleRepo := repository.NewRoleRepository(db)
		permRepo := repository.NewPermissionRepository(db)
		auditRepo := repository.NewAuditRepository(db)
		jwtMgr := infra.NewJWTTokenManager(jwtSecret, 2*time.Hour, 7*24*time.Hour, "erp-go")
		passHasher := infra.NewBcryptPasswordHasher(10)

		authService = app.NewAuthService(userRepo, roleRepo, jwtMgr, passHasher, auditRepo)
		userService = app.NewUserService(userRepo, roleRepo, passHasher, auditRepo)
		roleService = app.NewRoleService(roleRepo, permRepo, auditRepo)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(
		middleware.Recovery(log),
		middleware.RequestID(),
		middleware.TraceID(),
		middleware.TenantID(),
		middleware.UserID(),
		middleware.CORS(),
		middleware.RequestLogger(log),
	)

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": cfg.Server.Name,
			"db":      db != nil,
		})
	})

	if authService != nil {
		server := httpiface.NewServer(authService, userService, roleService, log, cfg, jwtSecret)
		server.RegisterRoutes(engine)
		log.Info("IAM 服务启动，数据库模式")
	} else {
		api := engine.Group("/api/v1/iam")
		{
			api.POST("/login", notImplYet)
			api.POST("/refresh", notImplYet)
			api.GET("/users", notImplYet)
			api.GET("/roles", notImplYet)
			api.GET("/permissions", notImplYet)
		}
		log.Info("IAM 服务启动，占位模式")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Infof("IAM 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("IAM 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭 IAM 服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("IAM 服务关闭异常: %v", err)
	}
	log.Info("IAM 服务已关闭")
}

func notImplYet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "接口已规划，数据库迁移完成后可用",
	})
}
