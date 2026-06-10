package repository

import (
	"github.com/Tangyd893/ERP-Go/backend/shared/config"
	"github.com/Tangyd893/ERP-Go/backend/shared/database"
	"gorm.io/gorm"
)

// NewDB 创建数据库连接（委托 shared/database）
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	return database.NewDB(cfg)
}
