package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "channel-service",
		DefaultPort: 8084,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			var channelAppService *app.ChannelAppService
			if db != nil {
				repo := repository.NewChannelRepository(db)
				channelAppService = app.NewChannelAppService(repo)
			}
			handler.NewChannelHandler(channelAppService).RegisterRoutes(engine.Group("/api/v1/channel"))
			return nil
		},
	})
	srv.WaitShutdown()
}
