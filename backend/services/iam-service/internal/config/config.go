package config

import (
	"time"

	sharedconfig "github.com/Tangyd893/ERP-Go/backend/shared/config"
)

// Config IAM 服务配置
type Config struct {
	Server   ServerConfig          `mapstructure:"server"`
	Database sharedconfig.DatabaseConfig `mapstructure:"database"`
	Redis    sharedconfig.RedisConfig    `mapstructure:"redis"`
	Log      sharedconfig.LogConfig      `mapstructure:"log"`
	JWT      JWTConfig             `mapstructure:"jwt"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	AccessExpiry     time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry    time.Duration `mapstructure:"refresh_expiry"`
	Issuer           string        `mapstructure:"issuer"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:         "iam-service",
			Host:         "0.0.0.0",
			Port:         8081,
			Mode:         "debug",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		JWT: JWTConfig{
			Secret:        "change-me-in-production",
			AccessExpiry:  2 * time.Hour,
			RefreshExpiry: 7 * 24 * time.Hour,
			Issuer:        "erp-go",
		},
	}
}
