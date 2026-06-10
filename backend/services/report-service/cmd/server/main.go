package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/report-service/internal/app"
	handler "github.com/Tangyd893/ERP-Go/backend/services/report-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "report-service",
		DefaultPort: 8093,
		RegisterRoutes: func(engine *gin.Engine, _ *gorm.DB, _ logger.Logger) error {
			reportAppService := app.NewReportAppService()
			handler.NewReportHandler(reportAppService).RegisterRoutes(engine.Group("/api/v1/report"))
			return nil
		},
	})
	srv.WaitShutdown()
}
