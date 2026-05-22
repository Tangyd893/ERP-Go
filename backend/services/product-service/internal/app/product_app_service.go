package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/infra/repository"
)

// ProductAppService 商品应用服务
type ProductAppService struct {
	repo *repository.ProductRepository
}

func NewProductAppService(repo *repository.ProductRepository) *ProductAppService {
	return &ProductAppService{repo: repo}
}

func (s *ProductAppService) CreateSKU(ctx context.Context, sku *domain.SKU) error {
	return s.repo.CreateSKU(ctx, sku)
}

func (s *ProductAppService) ListSKUs(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SKU, int64, error) {
	return s.repo.ListSKUs(ctx, tenantID, offset, limit)
}

func (s *ProductAppService) GetSKU(ctx context.Context, tenantID, code string) (*domain.SKU, error) {
	return s.repo.FindSKUByCode(ctx, tenantID, code)
}

func (s *ProductAppService) MapPlatformSKU(ctx context.Context, mapping *domain.PlatformSKU) error {
	return s.repo.MapPlatformSKU(ctx, mapping)
}

func (s *ProductAppService) GetPlatformSKU(ctx context.Context, tenantID, storeID, platformCode string) (*domain.PlatformSKU, error) {
	return s.repo.FindPlatformSKU(ctx, tenantID, storeID, platformCode)
}
