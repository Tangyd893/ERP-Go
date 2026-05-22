package outbox

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
)

// LogPublisher 日志记录型事件发布器（RabbitMQ 替换入口）
type LogPublisher struct {
	log logger.Logger
}

// NewLogPublisher 创建日志发布器，生产环境替换为 RabbitMQ 实现
func NewLogPublisher(log logger.Logger) *LogPublisher {
	return &LogPublisher{log: log}
}

func (p *LogPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	p.log.Infof("[事件发布] %s, payload=%s", eventType, string(payload))
	return nil
}

// Ensure LogPublisher implements EventPublisher
var _ EventPublisher = (*LogPublisher)(nil)
