package outbox

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
)

func TestMemOutboxStoreSave(t *testing.T) {
	store := NewMemOutboxStore()
	msg := &OutboxMessage{
		AggregateID: "order-1", AggregateType: "SalesOrder",
		TenantID: "default", EventType: events.EventOrderApproved,
		Payload: []byte(`{"test":true}`), Status: StatusPending,
		CreatedAt: time.Now(),
	}
	if err := store.Save(context.Background(), msg); err != nil {
		t.Fatalf("Save 失败: %v", err)
	}
	if msg.ID == 0 {
		t.Fatal("Save 应分配 ID")
	}
}

func TestMemOutboxStoreFetchPending(t *testing.T) {
	store := NewMemOutboxStore()
	saveMsg(t, store, StatusPending, "order-1")
	saveMsg(t, store, StatusPending, "order-2")
	saveMsg(t, store, StatusPublished, "order-3")

	pending, err := store.FetchPending(context.Background(), 10)
	if err != nil {
		t.Fatalf("FetchPending 失败: %v", err)
	}
	if len(pending) != 2 {
		t.Fatalf("应返回 2 条 pending，实际 %d", len(pending))
	}
}

func TestMemOutboxStoreMarkPublished(t *testing.T) {
	store := NewMemOutboxStore()
	msg := saveMsg(t, store, StatusPending, "order-1")

	if err := store.MarkPublished(context.Background(), msg.ID); err != nil {
		t.Fatalf("MarkPublished 失败: %v", err)
	}
	pending, _ := store.FetchPending(context.Background(), 10)
	if len(pending) != 0 {
		t.Fatal("已发布的消息不应在 pending 中")
	}
}

func TestMemOutboxStoreMarkFailed(t *testing.T) {
	store := NewMemOutboxStore()
	msg := saveMsg(t, store, StatusPending, "order-1")

	if err := store.MarkFailed(context.Background(), msg.ID, context.DeadlineExceeded); err != nil {
		t.Fatalf("MarkFailed 失败: %v", err)
	}
	failed, total, _ := store.FetchFailed(context.Background(), 0, 10)
	if total != 1 {
		t.Fatalf("应有 1 条 failed，实际 %d", total)
	}
	if failed[0].RetryCount != 1 {
		t.Fatalf("重试次数应为 1，实际 %d", failed[0].RetryCount)
	}
}

func TestMemOutboxStoreFetchFailed(t *testing.T) {
	store := NewMemOutboxStore()
	saveMsg(t, store, StatusFailed, "f-1")
	saveMsg(t, store, StatusFailed, "f-2")
	saveMsg(t, store, StatusFailed, "f-3")
	saveMsg(t, store, StatusPending, "p-1")

	failed, total, err := store.FetchFailed(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("FetchFailed 失败: %v", err)
	}
	if total != 3 {
		t.Fatalf("总数应为 3，实际 %d", total)
	}
	if len(failed) != 2 {
		t.Fatalf("分页应返回 2 条，实际 %d", len(failed))
	}
}

func TestMemInboxStoreIsDuplicate(t *testing.T) {
	store := NewMemInboxStore()
	ctx := context.Background()

	isDup, _ := store.IsDuplicate(ctx, "msg-1")
	if isDup {
		t.Fatal("新消息不应为重复")
	}
	store.Save(ctx, &InboxMessage{MessageID: "msg-1", EventType: events.EventOrderApproved, ProcessedAt: time.Now()})
	isDup, _ = store.IsDuplicate(ctx, "msg-1")
	if !isDup {
		t.Fatal("已保存消息应为重复")
	}
}

func TestMemInboxStoreSave(t *testing.T) {
	store := NewMemInboxStore()
	msg := &InboxMessage{
		MessageID:   "msg-inbox-1",
		EventType:   events.EventOutboundShipped,
		Payload:     []byte(`{"outbound_id":"OB-1"}`),
		ProcessedAt: time.Now(),
	}
	if err := store.Save(context.Background(), msg); err != nil {
		t.Fatalf("Save 失败: %v", err)
	}
}

func TestNewEventPayload(t *testing.T) {
	type testData struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}

	payload, err := NewEventPayload(events.EventOrderApproved, testData{
		OrderID: "order-1", Status: "approved",
	})
	if err != nil {
		t.Fatalf("NewEventPayload 失败: %v", err)
	}

	var ep EventPayload
	if err := json.Unmarshal(payload, &ep); err != nil {
		t.Fatalf("反序列化 EventPayload 失败: %v", err)
	}
	if ep.EventType != events.EventOrderApproved {
		t.Fatalf("事件类型应为 order.approved，实际 %s", ep.EventType)
	}

	var data testData
	if err := json.Unmarshal(ep.Data, &data); err != nil {
		t.Fatalf("反序列化 data 失败: %v", err)
	}
	if data.OrderID != "order-1" || data.Status != "approved" {
		t.Fatal("数据反序列化结果不正确")
	}
}

func TestEventPayloadRoundTrip(t *testing.T) {
	payload, err := NewEventPayload("test.event", map[string]int{"count": 42})
	if err != nil {
		t.Fatalf("NewEventPayload 失败: %v", err)
	}

	var ep EventPayload
	json.Unmarshal(payload, &ep)
	if ep.EventType != "test.event" {
		t.Fatalf("事件类型不匹配")
	}

	var data map[string]int
	json.Unmarshal(ep.Data, &data)
	if data["count"] != 42 {
		t.Fatalf("数据不匹配: %v", data)
	}
}

func TestOutboxProcessorProcessPending(t *testing.T) {
	store := NewMemOutboxStore()
	publisher := &logPublisher{}
	processor := NewOutboxProcessor(store, publisher, 5, 10*time.Millisecond)

	saveMsg(t, store, StatusPending, "order-1")
	saveMsg(t, store, StatusPending, "order-2")

	if err := processor.ProcessPending(context.Background()); err != nil {
		t.Fatalf("ProcessPending 失败: %v", err)
	}
	if publisher.count != 2 {
		t.Fatalf("应发布 2 条，实际 %d", publisher.count)
	}
}

func TestOutboxProcessorHandleInboxMessage(t *testing.T) {
	store := NewMemOutboxStore()
	inbox := NewMemInboxStore()
	publisher := &logPublisher{}
	processor := NewOutboxProcessor(store, publisher, 5, 10*time.Millisecond)

	handler := &testHandler{eventType: "order.approved"}
	processor.RegisterHandler(handler)

	payload, _ := NewEventPayload("order.approved", map[string]string{"order_id": "o-1"})

	err := processor.HandleInboxMessage(context.Background(), "msg-id-1", "order.approved", payload, inbox)
	if err != nil {
		t.Fatalf("HandleInboxMessage 失败: %v", err)
	}
	if !handler.called {
		t.Fatal("handler 应被调用")
	}

	handler.called = false
	err = processor.HandleInboxMessage(context.Background(), "msg-id-1", "order.approved", payload, inbox)
	if err != nil {
		t.Fatalf("重复消息应处理成功: %v", err)
	}
	if handler.called {
		t.Fatal("重复消息不应再次调用 handler")
	}
}

func TestOutboxMessageStatusConstants(t *testing.T) {
	if StatusPending != "pending" {
		t.Fatalf("StatusPending = %s", StatusPending)
	}
	if StatusPublished != "published" {
		t.Fatalf("StatusPublished = %s", StatusPublished)
	}
	if StatusFailed != "failed" {
		t.Fatalf("StatusFailed = %s", StatusFailed)
	}
}

func saveMsg(t *testing.T, store *MemOutboxStore, status MessageStatus, aggregateID string) *OutboxMessage {
	t.Helper()
	msg := &OutboxMessage{
		AggregateID:   aggregateID,
		AggregateType: "SalesOrder",
		TenantID:      "default",
		EventType:     events.EventOrderApproved,
		Payload:       []byte(`{}`),
		Status:        status,
		CreatedAt:     time.Now(),
	}
	if err := store.Save(context.Background(), msg); err != nil {
		t.Fatalf("saveMsg 失败: %v", err)
	}
	return msg
}

type logPublisher struct {
	count int
}

func (p *logPublisher) Publish(_ context.Context, eventType string, payload []byte) error {
	p.count++
	return nil
}

type testHandler struct {
	eventType string
	called    bool
}

func (h *testHandler) EventType() string { return h.eventType }
func (h *testHandler) Handle(_ context.Context, msg *OutboxMessage) error {
	h.called = true
	return nil
}
