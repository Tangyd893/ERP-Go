package outbox

import (
	"context"
	"sync"
)

// MemOutboxStore 内存实现的 OutboxStore，仅供测试使用
type MemOutboxStore struct {
	mu       sync.Mutex
	messages []*OutboxMessage
	nextID   int64
}

func NewMemOutboxStore() *MemOutboxStore {
	return &MemOutboxStore{messages: make([]*OutboxMessage, 0), nextID: 1}
}

func (s *MemOutboxStore) Save(ctx context.Context, msg *OutboxMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg.ID = s.nextID
	s.nextID++
	s.messages = append(s.messages, msg)
	return nil
}

func (s *MemOutboxStore) FetchPending(ctx context.Context, limit int) ([]*OutboxMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*OutboxMessage
	for _, msg := range s.messages {
		if msg.Status == StatusPending {
			result = append(result, msg)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (s *MemOutboxStore) MarkPublished(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, msg := range s.messages {
		if msg.ID == id {
			msg.Status = StatusPublished
			break
		}
	}
	return nil
}

func (s *MemOutboxStore) MarkFailed(ctx context.Context, id int64, err error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, msg := range s.messages {
		if msg.ID == id {
			msg.Status = StatusFailed
			msg.RetryCount++
			break
		}
	}
	return nil
}

func (s *MemOutboxStore) Retry(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, msg := range s.messages {
		if msg.ID == id {
			msg.Status = StatusPending
			break
		}
	}
	return nil
}

func (s *MemOutboxStore) FetchFailed(ctx context.Context, offset, limit int) ([]*OutboxMessage, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var failed []*OutboxMessage
	for _, msg := range s.messages {
		if msg.Status == StatusFailed {
			failed = append(failed, msg)
		}
	}
	total := int64(len(failed))
	if offset >= len(failed) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(failed) {
		end = len(failed)
	}
	return failed[offset:end], total, nil
}

// MemInboxStore 内存实现的 InboxStore，仅供测试使用
type MemInboxStore struct {
	mu       sync.Mutex
	messages map[string]bool
}

func NewMemInboxStore() *MemInboxStore {
	return &MemInboxStore{messages: make(map[string]bool)}
}

func (s *MemInboxStore) IsDuplicate(ctx context.Context, messageID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages[messageID], nil
}

func (s *MemInboxStore) Save(ctx context.Context, msg *InboxMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[msg.MessageID] = true
	return nil
}
