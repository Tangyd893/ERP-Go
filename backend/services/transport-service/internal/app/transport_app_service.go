package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/infra/repository"
)

type TransportAppService struct {
	repo *repository.TransportRepository
}

func NewTransportAppService(repo *repository.TransportRepository) *TransportAppService {
	return &TransportAppService{repo: repo}
}

func (s *TransportAppService) ListCarriers(ctx context.Context, tenantID string) ([]*domain.Carrier, error) {
	return s.repo.ListCarriers(ctx, tenantID)
}

func (s *TransportAppService) CreateShipment(ctx context.Context, s2 *domain.Shipment) error {
	return s.repo.CreateShipment(ctx, s2)
}

func (s *TransportAppService) ListShipments(ctx context.Context, tenantID string, offset, limit int) ([]*domain.Shipment, int64, error) {
	return s.repo.ListShipments(ctx, tenantID, offset, limit)
}

func (s *TransportAppService) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	return s.repo.FindShipment(ctx, id)
}

func (s *TransportAppService) CreateLabel(ctx context.Context, id, trackingNo, labelURL string) error {
	return s.repo.UpdateShipmentStatus(ctx, id, string(domain.ShipmentLabeled), trackingNo)
}
