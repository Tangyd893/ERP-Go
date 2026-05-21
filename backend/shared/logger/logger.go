package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// 标准日志字段
const (
	FieldServiceName = "service_name"
	FieldEnvironment = "environment"
	FieldTenantID    = "tenant_id"
	FieldUserID      = "user_id"
	FieldRequestID   = "request_id"
	FieldTraceID     = "trace_id"
	FieldSpanID      = "span_id"
	FieldErrorCode   = "error_code"
	FieldBusinessID  = "business_id"
	FieldDuration    = "duration_ms"
)

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
}

// Fields 日志字段类型
type Fields map[string]interface{}

// logrusLogger 基于 logrus 的日志实现
type logrusLogger struct {
	entry *logrus.Entry
}

// New 创建日志实例
func New(level, format, output, serviceName, environment string) Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FullTimestamp:   true,
		})
	}

	switch output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.SetOutput(os.Stdout)
			logger.Warnf("无法打开日志文件 %s，使用标准输出: %v", output, err)
		} else {
			logger.SetOutput(io.MultiWriter(os.Stdout, file))
		}
	}

	entry := logger.WithFields(logrus.Fields{
		FieldServiceName: serviceName,
		FieldEnvironment: environment,
	})

	return &logrusLogger{entry: entry}
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{entry: l.entry.WithField(key, value)}
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{entry: l.entry.WithError(err)}
}
