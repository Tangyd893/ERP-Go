package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "tenant-service",
		DefaultPort: 8082,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var tenantAppService *app.TenantAppService
			if db != nil {
				tenantRepo := repository.NewTenantRepository(db)
				orgRepo := repository.NewOrgRepository(db)
				tenantAppService = app.NewTenantAppService(tenantRepo, orgRepo)
			}
			handler.NewTenantHandler(tenantAppService).RegisterRoutes(engine.Group("/api/v1/tenant"))
			return nil
		},
	})
	srv.WaitShutdown()
}
