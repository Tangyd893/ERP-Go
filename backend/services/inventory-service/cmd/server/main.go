package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/inventory-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "inventory-service",
		DefaultPort: 8086,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var repo *repository.InventoryRepository
			if db != nil {
				repo = repository.NewInventoryRepository(db)
			}
			handler.NewInventoryHandler(repo).RegisterRoutes(engine.Group("/api/v1/inventory"))
			return nil
		},
	})
	srv.WaitShutdown()
}
