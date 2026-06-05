package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/warehouse-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	var appService *app.WarehouseAppService
	var coordinator *workflows.P4OutboundFlowCoordinator

	srv := server.New(server.Options{
		Name:        "warehouse-service",
		DefaultPort: 8087,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			db, err := repository.NewDB(cfg)
			if err != nil {
				return nil, err
			}
			appService = app.NewWarehouseAppService(repository.NewWarehouseRepository(db))
			orderURL := os.Getenv("ORDER_SERVICE_URL")
			if orderURL == "" {
				orderURL = "http://localhost:8085"
			}
			appService.WithFulfillmentClient(app.NewOrderFulfillmentClient(orderURL))
			return db, nil
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			wh := handler.NewWarehouseHandler(appService)
			if db != nil {
				outboxStore := outbox.NewPGOutboxStore(db)
				inboxStore := outbox.NewPGInboxStore(db)

				inventoryURL := os.Getenv("INVENTORY_SERVICE_URL")
				if inventoryURL == "" {
					inventoryURL = "http://localhost:8086"
				}
				inboundAdapter := workflows.NewHTTPInboundHandlerAdapter(inventoryURL)

				coordinator = workflows.NewP4OutboundFlowCoordinator(outboxStore, inboxStore)
				coordinator.SetInboundHandler(inboundAdapter)

				processor := outbox.NewOutboxProcessor(outboxStore, outbox.NewLogPublisher(log), 10, 5*time.Second)
				processor.RegisterHandler(&inboundReceivedHandler{coordinator: coordinator})

				ctx := context.Background()
				go outbox.StartPolling(ctx, processor, log)
				log.Info("Outbox 事件轮询已启动（采购入库）")

				wh.WithCoordinator(coordinator)
			}
			wh.RegisterRoutes(engine.Group("/api/v1/warehouse"))
			return nil
		},
	})
	srv.WaitShutdown()
}

type inboundReceivedHandler struct {
	coordinator *workflows.P4OutboundFlowCoordinator
}

func (h *inboundReceivedHandler) EventType() string { return "inventory.increased" }

func (h *inboundReceivedHandler) Handle(ctx context.Context, msg *outbox.OutboxMessage) error {
	return h.coordinator.HandleInboundReceived(ctx, strconv.FormatInt(msg.ID, 10), msg.Payload)
}
