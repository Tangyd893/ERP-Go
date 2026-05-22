package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "order-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8085
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	var db *gorm.DB
	var orderAppService *app.OrderAppService
	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用占位模式: %v", dbErr)
	} else {
		log.Info("数据库连接成功")
		db = database
		orderRepo := repository.NewOrderRepository(db)
		orderAppService = app.NewOrderAppService(orderRepo)

		outboxStore := outbox.NewPGOutboxStore(db)
		orderAppService.WithOutbox(outboxStore)

		publisher := outbox.NewLogPublisher(log)
		processor := outbox.NewOutboxProcessor(outboxStore, publisher, 10, 5*time.Second)

		inventoryURL := os.Getenv("INVENTORY_SERVICE_URL")
		if inventoryURL == "" {
			inventoryURL = "http://localhost:8086"
		}
		warehouseURL := os.Getenv("WAREHOUSE_SERVICE_URL")
		if warehouseURL == "" {
			warehouseURL = "http://localhost:8087"
		}

		stockAdapter := workflows.NewHTTPStockLockAdapter(inventoryURL)
		outboundAdapter := workflows.NewHTTPOutboundCreatorAdapter(warehouseURL)

		coordinator := workflows.NewP4OutboundFlowCoordinator(outboxStore, outbox.NewPGInboxStore(db))
		coordinator.SetStockHandler(stockAdapter)
		coordinator.SetOutboundCreator(outboundAdapter)

		processor.RegisterHandler(&orderApprovedHandler{coordinator: coordinator})
		processor.RegisterHandler(&orderCancelledHandler{coordinator: coordinator})

		ctx := context.Background()
		go outbox.StartPolling(ctx, processor, log)
		log.Info("Outbox 事件轮询已启动（订单审核→锁定库存→创建出库单）")
	}

	orderHandler := handler.NewOrderHandler(orderAppService)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(), middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log))

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{"status": status, "service": cfg.Server.Name, "db": db != nil})
	})

	orderHandler.RegisterRoutes(engine.Group("/api/v1/order"))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}

	go func() {
		log.Infof("Order 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Order 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭 Order 服务...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Order 服务关闭异常: %v", err)
	}
	log.Info("Order 服务已关闭")
}

// 事件处理器实现

type orderApprovedHandler struct {
	coordinator *workflows.P4OutboundFlowCoordinator
}

func (h *orderApprovedHandler) EventType() string { return "order.approved" }

func (h *orderApprovedHandler) Handle(ctx context.Context, msg *outbox.OutboxMessage) error {
	return h.coordinator.HandleOrderApproved(ctx, strconv.FormatInt(msg.ID, 10), msg.Payload)
}

type orderCancelledHandler struct {
	coordinator *workflows.P4OutboundFlowCoordinator
}

func (h *orderCancelledHandler) EventType() string { return "order.cancelled" }

func (h *orderCancelledHandler) Handle(ctx context.Context, msg *outbox.OutboxMessage) error {
	return h.coordinator.HandleOrderCancelled(ctx, strconv.FormatInt(msg.ID, 10), msg.Payload)
}
