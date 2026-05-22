package outbox

import (
	"context"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
)

// StartPolling 启动 Outbox 轮询 goroutine
// 定期从 outbox 表中取出 pending 消息，通过 publisher 发布
func StartPolling(ctx context.Context, processor *OutboxProcessor, log logger.Logger) {
	log.Info("Outbox 消息轮询已启动")

	ticker := time.NewTicker(processor.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Outbox 消息轮询已停止")
			return
		case <-ticker.C:
			if err := processor.ProcessPending(ctx); err != nil {
				log.Errorf("Outbox 消息处理失败: %v", err)
				continue
			}
			log.Debugf("Outbox 轮询: 批次处理完成")
		}
	}
}
