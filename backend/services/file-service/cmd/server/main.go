package main

import (
	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/infra/repository"
	handler "github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/interfaces/http"
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	srv := server.New(server.Options{
		Name:        "file-service",
		DefaultPort: 8089,
		InitDB: func(cfg config.DatabaseConfig, log logger.Logger) (*gorm.DB, error) {
			return repository.NewDB(cfg)
		},
		RegisterRoutes: func(engine *gin.Engine, db *gorm.DB, log logger.Logger) error {
			engine.MaxMultipartMemory = 32 << 20
			var fileAppService *app.FileAppService
			if db != nil {
				repo := repository.NewFileRepository(db)
				fileAppService = app.NewFileAppService(repo)
			}
			handler.NewFileHandler(fileAppService).RegisterRoutes(engine.Group("/api/v1/file"))
			return nil
		},
	})
	srv.WaitShutdown()
}
