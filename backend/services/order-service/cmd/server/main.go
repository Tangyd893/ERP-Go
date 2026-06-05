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
	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// orderDeps 持有 main 中初始化的全部依赖，便于在各阶段间传递
type orderDeps struct {
	db                  *gorm.DB
	srv                 *http.Server
	orderAppService     *app.OrderAppService
	orderHandler        *handler.OrderHandler
	publisher           outbox.EventPublisher
	coordinator         *workflows.P4OutboundFlowCoordinator
	fulfillmentConsumer *outbox.RabbitMQConsumer
}

func main() {
	cfg, _ := config.Load("")
	cfg.Server.Name = "order-service"
	if cfg.Server.Port == 0 || cfg.Server.Port == 8080 {
		cfg.Server.Port = 8085
	}
	log := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Server.Name, os.Getenv("ENVIRONMENT"))

	deps := setupDB(cfg, log)
	if deps.orderHandler == nil {
		deps.orderHandler = handler.NewOrderHandler(deps.orderAppService)
	}
	setupOutbox(cfg, deps, log)
	setupHTTP(cfg, log, deps)
	waitShutdown(log, deps)
}

// setupDB 初始化数据库连接与领域服务（仅在 DB 可用时）
func setupDB(cfg *config.Config, log logger.Logger) *orderDeps {
	deps := &orderDeps{}
	database, dbErr := repository.NewDB(cfg.Database)
	if dbErr != nil {
		log.Warnf("数据库连接失败，使用占位模式: %v", dbErr)
		return deps
	}
	log.Info("数据库连接成功")
	deps.db = database

	orderRepo := repository.NewOrderRepository(deps.db)
	deps.orderAppService = app.NewOrderAppService(orderRepo)
	deps.orderAppService.WithOutbox(outbox.NewPGOutboxStore(deps.db))
	return deps
}

// setupOutbox 初始化事件发布器、P4 编排器、Outbox 轮询及 RabbitMQ 消费者
func setupOutbox(cfg *config.Config, deps *orderDeps, log logger.Logger) {
	if deps.db == nil {
		return
	}
	outboxStore := outbox.NewPGOutboxStore(deps.db)

	rabbitURL := buildRabbitURL(cfg)
	publisher := newEventPublisher(rabbitURL, log)
	deps.publisher = publisher

	coordinator := buildCoordinator(deps.db, outboxStore, deps.orderAppService)
	deps.coordinator = coordinator

	processor := outbox.NewOutboxProcessor(outboxStore, publisher, 10, 5*time.Second)
	processor.RegisterHandler(&orderApprovedHandler{coordinator: coordinator})
	processor.RegisterHandler(&orderCancelledHandler{coordinator: coordinator})

	ctx := context.Background()
	go outbox.StartPolling(ctx, processor, log)
	log.Info("Outbox 事件轮询已启动（订单审核→锁定库存→创建出库单）")

	deps.fulfillmentConsumer = startFulfillmentConsumer(ctx, deps.db, rabbitURL, coordinator, log)

	deps.orderHandler = handler.NewOrderHandler(deps.orderAppService).
		WithCoordinator(coordinator).
		WithOutboxStore(outboxStore)
}

// setupHTTP 创建 Gin 引擎、注册中间件/路由、启动 HTTP 服务
func setupHTTP(cfg *config.Config, log logger.Logger, deps *orderDeps) {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(
		middleware.Recovery(log), middleware.RequestID(), middleware.TraceID(),
		middleware.TenantID(), middleware.CORS(), middleware.RequestLogger(log),
	)
	engine.GET("/health", healthHandler(cfg.Server.Name, deps.db))
	deps.orderHandler.RegisterRoutes(engine.Group("/api/v1/order"))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}

	go func() {
		log.Infof("Order 服务启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Order 服务启动失败: %v", err)
		}
	}()

	deps.srv = srv
}

// waitShutdown 等待信号并优雅关闭
func waitShutdown(log logger.Logger, deps *orderDeps) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭 Order 服务...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if deps.publisher != nil {
		if closer, ok := deps.publisher.(interface{ Close() error }); ok {
			closer.Close()
		}
	}
	if deps.fulfillmentConsumer != nil {
		deps.fulfillmentConsumer.Close()
	}
	if deps.db != nil {
		if sqlDB, err := deps.db.DB(); err == nil {
			sqlDB.Close()
		}
	}
	if deps.srv != nil {
		if err := deps.srv.Shutdown(shutdownCtx); err != nil {
			log.Errorf("Order 服务关闭异常: %v", err)
		}
	}
	log.Info("Order 服务已关闭")
}

// ---- 内部辅助函数 ----

func buildRabbitURL(cfg *config.Config) string {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" && cfg.RabbitMQ.Host != "" {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			cfg.RabbitMQ.User, cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Host, cfg.RabbitMQ.Port, cfg.RabbitMQ.VHost)
	}
	return rabbitURL
}

func newEventPublisher(rabbitURL string, log logger.Logger) outbox.EventPublisher {
	if rabbitURL == "" {
		log.Info("未配置 RabbitMQ，使用日志发布模式")
		return outbox.NewLogPublisher(log)
	}
	rmqPublisher, rmqErr := outbox.NewRabbitMQPublisher(rabbitURL, "erp.events", log)
	if rmqErr != nil {
		log.Warnf("RabbitMQ 发布器初始化失败，降级为日志发布: %v", rmqErr)
		return outbox.NewLogPublisher(log)
	}
	log.Info("RabbitMQ 事件发布器已就绪")
	return rmqPublisher
}

func buildCoordinator(db *gorm.DB, outboxStore outbox.OutboxStore, orderAppService *app.OrderAppService) *workflows.P4OutboundFlowCoordinator {
	inventoryURL := os.Getenv("INVENTORY_SERVICE_URL")
	if inventoryURL == "" {
		inventoryURL = "http://localhost:8086"
	}
	warehouseURL := os.Getenv("WAREHOUSE_SERVICE_URL")
	if warehouseURL == "" {
		warehouseURL = "http://localhost:8087"
	}

	orderRepo := repository.NewOrderRepository(db)
	coordinator := workflows.NewP4OutboundFlowCoordinator(outboxStore, outbox.NewPGInboxStore(db))
	coordinator.SetStockHandler(workflows.NewHTTPStockLockAdapter(inventoryURL))
	coordinator.SetStockDeductHandler(workflows.NewHTTPStockDeductAdapter(inventoryURL))
	coordinator.SetOutboundCreator(workflows.NewHTTPOutboundCreatorAdapter(warehouseURL))
	coordinator.SetOrderStatusUpdater(app.NewLocalOrderStatusUpdater(orderRepo))
	return coordinator
}

func startFulfillmentConsumer(ctx context.Context, db *gorm.DB, rabbitURL string,
	coordinator *workflows.P4OutboundFlowCoordinator, log logger.Logger) *outbox.RabbitMQConsumer {
	if rabbitURL == "" {
		return nil
	}
	inboxStore := outbox.NewPGInboxStore(db)
	consumer, err := outbox.NewRabbitMQConsumer(
		ctx, rabbitURL, "order.fulfillment",
		[]string{events.EventOutboundShipped},
		func(cctx context.Context, eventType, messageID string, payload []byte) error {
			if eventType != events.EventOutboundShipped {
				return nil
			}
			return coordinator.HandleOutboundShipped(cctx, messageID, payload)
		},
		inboxStore, log,
	)
	if err != nil {
		log.Warnf("RabbitMQ 履约消费者启动失败: %v", err)
		return nil
	}
	log.Info("RabbitMQ 履约消费者已启动: queue=order.fulfillment")
	return consumer
}

func healthHandler(name string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{"status": status, "service": name, "db": db != nil})
	}
}

// ---- 事件处理器实现 ----

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
