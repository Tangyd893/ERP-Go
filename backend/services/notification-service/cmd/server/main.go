package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/notification-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "notification-service",
		DefaultPort: 8094,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var notificationAppService *app.NotificationAppService
			if db != nil {
				repo := repository.NewNotificationRepository(db)
				notificationAppService = app.NewNotificationAppService(repo)
			}
			handler.NewNotificationHandler(notificationAppService).RegisterRoutes(engine.Group("/api/v1/notification"))
			return nil
		},
	})
	srv.WaitShutdown()
}
