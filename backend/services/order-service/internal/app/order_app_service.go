package app

import (
	"context"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	sharedEvents "github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
)

// orderRepo 订单仓储接口
type orderRepo interface {
	Create(ctx context.Context, order *domain.SalesOrder) error
	FindByID(ctx context.Context, id string) (*domain.SalesOrder, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SalesOrder, int64, error)
}

// OrderAppService 订单应用服务
type OrderAppService struct {
	repo   orderRepo
	outbox outbox.OutboxStore
}

func NewOrderAppService(repo orderRepo) *OrderAppService {
	return &OrderAppService{repo: repo}
}

// WithOutbox 注入 Outbox 存储，启用事件发布
func (s *OrderAppService) WithOutbox(store outbox.OutboxStore) *OrderAppService {
	s.outbox = store
	return s
}

func (s *OrderAppService) ListOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SalesOrder, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}

func (s *OrderAppService) GetOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderAppService) CreateOrder(ctx context.Context, order *domain.SalesOrder) error {
	if err := s.repo.Create(ctx, order); err != nil {
		return err
	}
	s.emitEvent(ctx, order.ID, order.TenantID, "SalesOrder", sharedEvents.EventOrderImported, order)
	return nil
}

func (s *OrderAppService) ApproveOrder(ctx context.Context, id, operator string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.Approve(operator); err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(ctx, id, string(order.Status)); err != nil {
		return err
	}
	s.emitEvent(ctx, order.ID, order.TenantID, "SalesOrder", sharedEvents.EventOrderApproved, order)
	return nil
}

func (s *OrderAppService) CancelOrder(ctx context.Context, id, operator, reason string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.Cancel(operator, reason); err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(ctx, id, string(order.Status)); err != nil {
		return err
	}
	s.emitEvent(ctx, order.ID, order.TenantID, "SalesOrder", sharedEvents.EventOrderCancelled, order)
	return nil
}

func (s *OrderAppService) MarkAbnormal(ctx context.Context, id, operator, reason string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.MarkAbnormal(operator, reason); err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(ctx, id, string(order.Status)); err != nil {
		return err
	}
	s.emitEvent(ctx, order.ID, order.TenantID, "SalesOrder", sharedEvents.EventOrderAbnormal, order)
	return nil
}

func (s *OrderAppService) emitEvent(ctx context.Context, aggregateID, tenantID, aggregateType, eventType string, data interface{}) {
	if s.outbox == nil {
		return
	}
	payload, err := outbox.NewEventPayload(eventType, data)
	if err != nil {
		return
	}
	_ = s.outbox.Save(ctx, &outbox.OutboxMessage{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		TenantID:      tenantID,
		EventType:     eventType,
		Payload:       payload,
		Status:        outbox.StatusPending,
		CreatedAt:     time.Now(),
	})
}
