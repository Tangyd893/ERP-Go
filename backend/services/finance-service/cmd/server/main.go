package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/finance-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "finance-service",
		DefaultPort: 8092,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var financeAppService *app.FinanceAppService
			if db != nil {
				repo := repository.NewFinanceRepository(db)
				financeAppService = app.NewFinanceAppService(repo)
			}
			handler.NewFinanceHandler(financeAppService).RegisterRoutes(engine.Group("/api/v1/finance"))
			return nil
		},
	})
	srv.WaitShutdown()
}
