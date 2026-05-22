package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/infra/repository"
)

// OrderAppService 订单应用服务
type OrderAppService struct {
	repo *repository.OrderRepository
}

func NewOrderAppService(repo *repository.OrderRepository) *OrderAppService {
	return &OrderAppService{repo: repo}
}

func (s *OrderAppService) ListOrders(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SalesOrder, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}

func (s *OrderAppService) GetOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderAppService) CreateOrder(ctx context.Context, order *domain.SalesOrder) error {
	return s.repo.Create(ctx, order)
}

func (s *OrderAppService) ApproveOrder(ctx context.Context, id, operator string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.Approve(operator); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, string(order.Status))
}

func (s *OrderAppService) CancelOrder(ctx context.Context, id, operator, reason string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.Cancel(operator, reason); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, string(order.Status))
}

func (s *OrderAppService) MarkAbnormal(ctx context.Context, id, operator, reason string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := order.MarkAbnormal(operator, reason); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, string(order.Status))
}
