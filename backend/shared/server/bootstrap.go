package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Options 服务启动配置
type Options struct {
	// 服务名称（必填）
	Name string
	// 默认端口（当 config 为 0 或 8080 时使用此值）
	DefaultPort int
	// 数据库初始化函数，返回 nil 表示"不连DB"或"连接失败"
	InitDB func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error)
	// 路由注册函数，engine 和 db（可能为 nil）和 log
	RegisterRoutes func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error
	// 关闭前清理回调（如关闭 publisher、consumer 等）
	OnShutdown func()
	// HTTP 读超时
	ReadTimeout time.Duration
	// HTTP 写超时
	WriteTimeout time.Duration
}

// Server 封装了 Gin engine + HTTP server + 优雅关闭
type Server struct {
	engine *gin.Engine
	srv    *http.Server
	log    logger.Logger
	db     *gorm.DB
	opts   Options
}

// New 创建并启动服务
func New(opts Options) *Server {
	cfg, _ := config.Load("")
	cfg.Server.Name = opts.Name
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = opts.DefaultPort
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.Recovery(log),
		middleware.RequestID(),
		middleware.TraceID(),
		middleware.TenantID(),
		middleware.CORS(),
		middleware.RequestLogger(log),
	)

	var db *gorm.DB
	if opts.InitDB != nil {
		database, err := opts.InitDB(cfg.Database, log)
		if err != nil {
			log.Warnf("数据库连接失败，使用占位模式: %v", err)
		} else {
			log.Info("数据库连接成功")
			db = database
		}
	}

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": opts.Name,
			"db":      db != nil,
		})
	})

	if opts.RegisterRoutes != nil {
		if err := opts.RegisterRoutes(engine, db, log); err != nil {
			log.Fatalf("路由注册失败: %v", err)
		}
	}

	readTimeout := opts.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 30 * time.Second
	}
	writeTimeout := opts.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 30 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		log.Infof("%s 启动在 %s", opts.Name, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s 启动失败: %v", opts.Name, err)
		}
	}()

	return &Server{engine: engine, srv: srv, log: log, db: db, opts: opts}
}

// WaitShutdown 等待信号并执行优雅关闭
func (s *Server) WaitShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.log.Infof("正在关闭 %s...", s.opts.Name)
	if s.opts.OnShutdown != nil {
		s.opts.OnShutdown()
	}
	if s.db != nil {
		if sqlDB, _ := s.db.DB(); sqlDB != nil {
			sqlDB.Close()
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Errorf("%s 关闭异常: %v", s.opts.Name, err)
	}
	s.log.Infof("%s 已关闭", s.opts.Name)
}

// Engine 返回 gin engine
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// DB 返回数据库连接
func (s *Server) DB() *gorm.DB {
	return s.db
}

// Log 返回日志器
func (s *Server) Log() logger.Logger {
	return s.log
}
