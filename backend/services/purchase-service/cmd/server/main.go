package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/purchase-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "purchase-service",
		DefaultPort: 8091,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var purchaseAppService *app.PurchaseAppService
			if db != nil {
				repo := repository.NewPurchaseRepository(db)
				purchaseAppService = app.NewPurchaseAppService(repo)
			}
			handler.NewPurchaseHandler(purchaseAppService).RegisterRoutes(engine.Group("/api/v1/purchase"))
			return nil
		},
	})
	srv.WaitShutdown()
}
