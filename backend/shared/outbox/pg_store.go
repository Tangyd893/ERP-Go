package outbox

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const whereStatus = "status = ?"

// OutboxMessageModel Outbox 消息 GORM 模型
type OutboxMessageModel struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	AggregateID   string     `gorm:"column:aggregate_id;index"`
	AggregateType string     `gorm:"column:aggregate_type"`
	TenantID      string     `gorm:"column:tenant_id;index"`
	EventType     string     `gorm:"column:event_type;index"`
	Payload       []byte     `gorm:"column:payload"`
	Status        string     `gorm:"column:status;index"`
	RetryCount    int        `gorm:"column:retry_count"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	PublishedAt   *time.Time `gorm:"column:published_at"`
}

func (OutboxMessageModel) TableName() string { return "outbox_messages" }

// InboxMessageModel Inbox 消息 GORM 模型
type InboxMessageModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	MessageID   string    `gorm:"column:message_id;uniqueIndex"`
	EventType   string    `gorm:"column:event_type"`
	Payload     []byte    `gorm:"column:payload"`
	ProcessedAt time.Time `gorm:"column:processed_at"`
}

func (InboxMessageModel) TableName() string { return "inbox_messages" }

// PGOutboxStore PostgreSQL Outbox 存储实现
type PGOutboxStore struct {
	db *gorm.DB
}

func NewPGOutboxStore(db *gorm.DB) *PGOutboxStore {
	return &PGOutboxStore{db: db}
}

func (s *PGOutboxStore) Save(ctx context.Context, msg *OutboxMessage) error {
	return s.db.WithContext(ctx).Create(&OutboxMessageModel{
		AggregateID:   msg.AggregateID,
		AggregateType: msg.AggregateType,
		TenantID:      msg.TenantID,
		EventType:     msg.EventType,
		Payload:       msg.Payload,
		Status:        string(msg.Status),
		RetryCount:    msg.RetryCount,
		CreatedAt:     msg.CreatedAt,
	}).Error
}

func (s *PGOutboxStore) FetchPending(ctx context.Context, limit int) ([]*OutboxMessage, error) {
	var models []*OutboxMessageModel
	err := s.db.WithContext(ctx).
		Where("status = ?", string(StatusPending)).
		Order("created_at ASC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	messages := make([]*OutboxMessage, len(models))
	for i, m := range models {
		messages[i] = &OutboxMessage{
			ID:            m.ID,
			AggregateID:   m.AggregateID,
			AggregateType: m.AggregateType,
			TenantID:      m.TenantID,
			EventType:     m.EventType,
			Payload:       m.Payload,
			Status:        MessageStatus(m.Status),
			RetryCount:    m.RetryCount,
			CreatedAt:     m.CreatedAt,
			PublishedAt:   m.PublishedAt,
		}
	}
	return messages, nil
}

func (s *PGOutboxStore) MarkPublished(ctx context.Context, id int64) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&OutboxMessageModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       string(StatusPublished),
		"published_at": &now,
	}).Error
}

func (s *PGOutboxStore) MarkFailed(ctx context.Context, id int64, err error) error {
	return s.db.WithContext(ctx).Model(&OutboxMessageModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      string(StatusFailed),
		"retry_count": gorm.Expr("retry_count + 1"),
	}).Error
}

func (s *PGOutboxStore) Retry(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Model(&OutboxMessageModel{}).Where("id = ?", id).Update("status", string(StatusPending)).Error
}

func (s *PGOutboxStore) FetchFailed(ctx context.Context, offset, limit int) ([]*OutboxMessage, int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&OutboxMessageModel{}).
		Where("status = ?", string(StatusFailed)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*OutboxMessageModel
	err := s.db.WithContext(ctx).
		Where("status = ?", string(StatusFailed)).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	messages := make([]*OutboxMessage, len(models))
	for i, m := range models {
		messages[i] = &OutboxMessage{
			ID:            m.ID,
			AggregateID:   m.AggregateID,
			AggregateType: m.AggregateType,
			TenantID:      m.TenantID,
			EventType:     m.EventType,
			Payload:       m.Payload,
			Status:        MessageStatus(m.Status),
			RetryCount:    m.RetryCount,
			CreatedAt:     m.CreatedAt,
			PublishedAt:   m.PublishedAt,
		}
	}
	return messages, total, nil
}

// PGInboxStore PostgreSQL Inbox 存储实现
type PGInboxStore struct {
	db *gorm.DB
}

func NewPGInboxStore(db *gorm.DB) *PGInboxStore {
	return &PGInboxStore{db: db}
}

func (s *PGInboxStore) IsDuplicate(ctx context.Context, messageID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&InboxMessageModel{}).Where("message_id = ?", messageID).Count(&count).Error
	return count > 0, err
}

func (s *PGInboxStore) Save(ctx context.Context, msg *InboxMessage) error {
	return s.db.WithContext(ctx).Create(&InboxMessageModel{
		MessageID:   msg.MessageID,
		EventType:   msg.EventType,
		Payload:     msg.Payload,
		ProcessedAt: msg.ProcessedAt,
	}).Error
}

// EventPayload 事件载荷通用结构
type EventPayload struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

// NewEventPayload 构建事件载荷
func NewEventPayload(eventType string, data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&EventPayload{
		EventType: eventType,
		Data:      dataBytes,
	})
}
