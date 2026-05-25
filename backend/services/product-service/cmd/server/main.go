package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "product-service",
		DefaultPort: 8083,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var productAppService *app.ProductAppService
			if db != nil {
				repo := repository.NewProductRepository(db)
				productAppService = app.NewProductAppService(repo)
			}
			handler.NewProductHandler(productAppService).RegisterRoutes(engine.Group("/api/v1/product"))
			return nil
		},
	})
	srv.WaitShutdown()
}
