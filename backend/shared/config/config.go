package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// IsProduction 返回当前是否生产环境
func IsProduction() bool {
	return os.Getenv("ENVIRONMENT") == "production"
}

// IsDevelopment 返回当前是否开发环境（ENVIRONMENT 为空或 "development"）
func IsDevelopment() bool {
	env := os.Getenv("ENVIRONMENT")
	return env == "" || env == "development"
}

// ValidateProduction 生产环境关键配置校验，返回 nil 表示通过
func ValidateProduction(jwtSecret, dbPassword string) error {
	errs := make([]string, 0)

	if jwtSecret == "" || jwtSecret == DefaultJWTSecret || jwtSecret == DefaultJWTSecretLegacy {
		errs = append(errs, "JWT_SECRET 为默认值，生产环境必须设置唯一的 JWT 密钥")
	}

	if dbPassword == "" || dbPassword == DefaultDBPassword {
		errs = append(errs, "DATABASE_PASSWORD 为默认值，生产环境必须设置数据库密码")
	}

	if len(errs) > 0 {
		return fmt.Errorf("生产环境安全校验失败:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

// DefaultJWTSecret IAM 默认 JWT Secret（与 gateway 的默认值统一）
const DefaultJWTSecret = "erp-go-dev-secret-change-in-production"

// DefaultJWTSecretLegacy gateway 旧版默认值（兼容校验）
const DefaultJWTSecretLegacy = "erp-go-default-secret-change-in-production"

// DefaultDBPassword 默认数据库密码（仅开发环境可用）
const DefaultDBPassword = "erp123"

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

	// T-609: 生产环境数据库密码不能是默认值
	if IsProduction() && cfg.Database.Password == DefaultDBPassword {
		return nil, fmt.Errorf("生产环境安全校验失败: DATABASE_PASSWORD 为默认值 %q，必须设置数据库密码", DefaultDBPassword)
	}

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	setInt("SERVER_PORT", &cfg.Server.Port)
	setString("SERVER_HOST", &cfg.Server.Host)
	setString("SERVER_MODE", &cfg.Server.Mode)
	setString("DATABASE_HOST", &cfg.Database.Host)
	setInt("DATABASE_PORT", &cfg.Database.Port)
	setString("DATABASE_USER", &cfg.Database.User)
	setString("DATABASE_PASSWORD", &cfg.Database.Password)
	setString("DATABASE_DBNAME", &cfg.Database.DBName)
	setString("DATABASE_SSLMODE", &cfg.Database.SSLMode)
	setString("REDIS_HOST", &cfg.Redis.Host)
	setInt("REDIS_PORT", &cfg.Redis.Port)
	setString("REDIS_PASSWORD", &cfg.Redis.Password)
	setString("RABBITMQ_HOST", &cfg.RabbitMQ.Host)
	setInt("RABBITMQ_PORT", &cfg.RabbitMQ.Port)
	setString("LOG_LEVEL", &cfg.Log.Level)
	setString("LOG_FORMAT", &cfg.Log.Format)
}

func setString(envKey string, target *string) {
	if v := os.Getenv(envKey); v != "" {
		*target = v
	}
}

func setInt(envKey string, target *int) {
	if v := os.Getenv(envKey); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			*target = port
		}
	}
}
