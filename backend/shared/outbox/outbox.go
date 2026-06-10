package outbox

import (
	"context"
	"time"
)

// MessageStatus 消息状态
type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusPublished MessageStatus = "published"
	StatusFailed    MessageStatus = "failed"
)

// OutboxMessage 发件箱消息
type OutboxMessage struct {
	ID            int64         `json:"id"`
	AggregateID   string        `json:"aggregate_id"`
	AggregateType string        `json:"aggregate_type"`
	TenantID      string        `json:"tenant_id"`
	EventType     string        `json:"event_type"`
	Payload       []byte        `json:"payload"`
	Status        MessageStatus `json:"status"`
	RetryCount    int           `json:"retry_count"`
	CreatedAt     time.Time     `json:"created_at"`
	PublishedAt   *time.Time    `json:"published_at,omitempty"`
}

// InboxMessage 收件箱消息（幂等消费）
type InboxMessage struct {
	ID          int64     `json:"id"`
	MessageID   string    `json:"message_id"`
	EventType   string    `json:"event_type"`
	Payload     []byte    `json:"payload"`
	ProcessedAt time.Time `json:"processed_at"`
}

// OutboxStore Outbox 存储接口
type OutboxStore interface {
	Save(ctx context.Context, msg *OutboxMessage) error
	FetchPending(ctx context.Context, limit int) ([]*OutboxMessage, error)
	FetchFailed(ctx context.Context, offset, limit int) ([]*OutboxMessage, int64, error)
	MarkPublished(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64, err error) error
	Retry(ctx context.Context, id int64) error
}

// InboxStore Inbox 存储接口
type InboxStore interface {
	IsDuplicate(ctx context.Context, messageID string) (bool, error)
	Save(ctx context.Context, msg *InboxMessage) error
}

// EventPublisher 事件发布器接口
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, msg *OutboxMessage) error
	EventType() string
}

// OutboxProcessor Outbox 消息处理调度器
type OutboxProcessor struct {
	store      OutboxStore
	publisher  EventPublisher
	handlers   map[string]EventHandler
	batchSize  int
	pollInterval time.Duration
}

// NewOutboxProcessor 创建 Outbox 处理器
func NewOutboxProcessor(store OutboxStore, publisher EventPublisher, batchSize int, pollInterval time.Duration) *OutboxProcessor {
	return &OutboxProcessor{
		store:        store,
		publisher:    publisher,
		handlers:     make(map[string]EventHandler),
		batchSize:    batchSize,
		pollInterval: pollInterval,
	}
}

// RegisterHandler 注册事件处理器
func (p *OutboxProcessor) RegisterHandler(handler EventHandler) {
	p.handlers[handler.EventType()] = handler
}

// ProcessPending 处理待发送消息
func (p *OutboxProcessor) ProcessPending(ctx context.Context) error {
	messages, err := p.store.FetchPending(ctx, p.batchSize)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err := p.publishMessage(ctx, msg); err != nil {
			if markErr := p.store.MarkFailed(ctx, msg.ID, err); markErr != nil {
				return markErr
			}
			continue
		}
		if err := p.store.MarkPublished(ctx, msg.ID); err != nil {
			return err
		}
	}

	return nil
}

// HandleInboxMessage 处理收件箱消息（幂等消费）
func (p *OutboxProcessor) HandleInboxMessage(ctx context.Context, messageID, eventType string, payload []byte, inboxStore InboxStore) error {
	isDuplicate, err := inboxStore.IsDuplicate(ctx, messageID)
	if err != nil {
		return err
	}
	if isDuplicate {
		return nil
	}

	handler, ok := p.handlers[eventType]
	if !ok {
		return nil
	}

	if err := handler.Handle(ctx, &OutboxMessage{
		EventType: eventType,
		Payload:   payload,
	}); err != nil {
		return err
	}

	return inboxStore.Save(ctx, &InboxMessage{
		MessageID:   messageID,
		EventType:   eventType,
		Payload:     payload,
		ProcessedAt: time.Now(),
	})
}

func (p *OutboxProcessor) publishMessage(ctx context.Context, msg *OutboxMessage) error {
	return p.publisher.Publish(ctx, msg.EventType, msg.Payload)
}
