package main

import (
	"os"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "transport-service",
		DefaultPort: 8088,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var transportAppService *app.TransportAppService
			if db != nil {
				repo := repository.NewTransportRepository(db)
				transportAppService = app.NewTransportAppService(repo)

				channelURL := os.Getenv("CHANNEL_SERVICE_URL")
				if channelURL == "" {
					channelURL = "http://localhost:8082"
				}
				transportAppService.WithChannelClient(app.NewChannelNotifyClient(channelURL))
			}
			handler.NewTransportHandler(transportAppService).RegisterRoutes(engine.Group("/api/v1/transport"))
			return nil
		},
	})
	srv.WaitShutdown()
}
