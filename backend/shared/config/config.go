package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 基础服务配置，所有服务共享
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务基础配置
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"db_name"`
	Schema          string `mapstructure:"schema"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// DSN 返回数据库连接字符串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// Addr 返回 Redis 地址
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// URL 返回 RabbitMQ 连接 URL
func (r *RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.User, r.Password, r.Host, r.Port, r.VHost)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:         "erp-service",
			Host:         "0.0.0.0",
			Port:         8080,
			Mode:         "debug",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "erp",
			Password:        "erp123",
			DBName:          "erp_go",
			Schema:          "public",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 300,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			PoolSize: 10,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     "localhost",
			Port:     5672,
			User:     "guest",
			Password: "guest",
			VHost:    "/",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}

// Load 从文件加载配置
func Load(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("DATABASE_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DATABASE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DATABASE_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DATABASE_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DATABASE_DBNAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("DATABASE_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Redis.Port = port
		}
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("RABBITMQ_HOST"); v != "" {
		cfg.RabbitMQ.Host = v
	}
	if v := os.Getenv("RABBITMQ_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.RabbitMQ.Port = port
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Log.Format = v
	}
}
