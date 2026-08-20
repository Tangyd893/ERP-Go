package workflows

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

func TestNewEventPayload(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	data := testData{Name: "test", Value: 42}

	payload, err := outbox.NewEventPayload("test.event", data)
	if err != nil {
		t.Fatalf("构建事件载荷失败: %v", err)
	}

	var ep outbox.EventPayload
	if err := json.Unmarshal(payload, &ep); err != nil {
		t.Fatalf("解析事件载荷失败: %v", err)
	}
	if ep.EventType != "test.event" {
		t.Errorf("事件类型应为 test.event，实际: %s", ep.EventType)
	}

	var parsed testData
	if err := json.Unmarshal(ep.Data, &parsed); err != nil {
		t.Fatalf("解析数据载荷失败: %v", err)
	}
	if parsed.Value != 42 {
		t.Errorf("值应为 42，实际: %d", parsed.Value)
	}
}

func TestMemOutboxStore(t *testing.T) {
	store := outbox.NewMemOutboxStore()

	msg := &outbox.OutboxMessage{
		AggregateID:   "order-001",
		AggregateType: "SalesOrder",
		EventType:     "order.created",
		Payload:       []byte(`{"test":true}`),
		Status:        outbox.StatusPending,
		CreatedAt:     time.Now(),
	}
	if err := store.Save(context.Background(), msg); err != nil {
		t.Fatalf("保存消息失败: %v", err)
	}
	if msg.ID != 1 {
		t.Errorf("消息 ID 应为 1，实际: %d", msg.ID)
	}

	pending, err := store.FetchPending(context.Background(), 10)
	if err != nil {
		t.Fatalf("获取 pending 消息失败: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("应有 1 条 pending 消息，实际: %d", len(pending))
	}

	if err := store.MarkPublished(context.Background(), msg.ID); err != nil {
		t.Fatalf("标记已发布失败: %v", err)
	}

	pending, _ = store.FetchPending(context.Background(), 10)
	if len(pending) != 0 {
		t.Errorf("标记已发布后 pending 消息应为 0，实际: %d", len(pending))
	}
}

func TestMemInboxStore(t *testing.T) {
	store := outbox.NewMemInboxStore()

	dup, err := store.IsDuplicate(context.Background(), "msg-001")
	if err != nil {
		t.Fatalf("检查重复失败: %v", err)
	}
	if dup {
		t.Error("新消息不应为重复")
	}

	if err := store.Save(context.Background(), &outbox.InboxMessage{
		MessageID: "msg-001", EventType: "test.event",
		Payload: []byte(`{}`), ProcessedAt: time.Now(),
	}); err != nil {
		t.Fatalf("保存 inbox 消息失败: %v", err)
	}

	dup, _ = store.IsDuplicate(context.Background(), "msg-001")
	if !dup {
		t.Error("已保存的消息应标记为重复")
	}
}
