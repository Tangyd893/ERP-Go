package repository

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/product-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateSPU(ctx context.Context, spu *domain.SPU) error {
	imagesJSON := ""
	if len(spu.Images) > 0 {
		imagesJSON = spu.Images[0]
	}
	return r.db.WithContext(ctx).Create(&SPUModel{
		ID: spu.ID, TenantID: spu.TenantID, Name: spu.Name,
		CategoryID: spu.CategoryID, Brand: spu.Brand, Description: spu.Description,
		MainImage: imagesJSON, Status: spu.Status, CreatedAt: spu.CreatedAt, UpdatedAt: spu.UpdatedAt,
	}).Error
}

func (r *ProductRepository) CreateSKU(ctx context.Context, sku *domain.SKU) error {
	return r.db.WithContext(ctx).Create(&SKUModel{
		ID: sku.ID, TenantID: sku.TenantID, SPUID: sku.SPUID,
		SKUCode: sku.Code, Barcode: sku.Barcode, SpecDesc: sku.SpecDesc,
		WeightGram: int(sku.Weight), LengthCm: sku.Length, WidthCm: sku.Width, HeightCm: sku.Height,
		PurchasePrice: sku.PurchasePrice, SellingPrice: sku.SalePrice, Currency: sku.Currency,
		Status: sku.Status, CreatedAt: sku.CreatedAt, UpdatedAt: sku.UpdatedAt,
	}).Error
}

func (r *ProductRepository) FindSKUByCode(ctx context.Context, tenantID, code string) (*domain.SKU, error) {
	var m SKUModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND sku_code = ?", tenantID, code).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewBusinessError(errors.CodeSKUNotFound, "SKU不存在")
		}
		return nil, err
	}
	return &domain.SKU{
		ID: m.ID, TenantID: m.TenantID, SPUID: m.SPUID, Code: m.SKUCode, Name: m.SKUCode,
		Barcode: m.Barcode, SpecDesc: m.SpecDesc,
		Weight: float64(m.WeightGram), Length: m.LengthCm, Width: m.WidthCm, Height: m.HeightCm,
		PurchasePrice: m.PurchasePrice, SalePrice: m.SellingPrice, Currency: m.Currency,
		Status: m.Status, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *ProductRepository) ListSKUs(ctx context.Context, tenantID string, offset, limit int) ([]*domain.SKU, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&SKUModel{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	var models []*SKUModel
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models)
	skus := make([]*domain.SKU, len(models))
	for i, m := range models {
		skus[i] = &domain.SKU{
			ID: m.ID, TenantID: m.TenantID, SPUID: m.SPUID, Code: m.SKUCode, Name: m.SKUCode,
			Barcode: m.Barcode, SpecDesc: m.SpecDesc,
			Weight: float64(m.WeightGram), Length: m.LengthCm, Width: m.WidthCm, Height: m.HeightCm,
			PurchasePrice: m.PurchasePrice, SalePrice: m.SellingPrice, Currency: m.Currency,
			Status: m.Status, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		}
	}
	return skus, total, nil
}

func (r *ProductRepository) MapPlatformSKU(ctx context.Context, mapping *domain.PlatformSKU) error {
	return r.db.WithContext(ctx).Create(&PlatformSKUModel{
		ID: mapping.ID, TenantID: mapping.TenantID, SKUID: mapping.SKUID,
		StoreID: mapping.StoreID, PlatformCode: mapping.PlatformCode,
		PlatformSKUID: mapping.PlatformSKU, ASIN: mapping.ASIN, FNSKU: mapping.FNSKU,
		PlatformStatus: mapping.PlatformStatus, CreatedAt: mapping.CreatedAt, UpdatedAt: mapping.UpdatedAt,
	}).Error
}

func (r *ProductRepository) FindPlatformSKU(ctx context.Context, tenantID, storeID, platformCode string) (*domain.PlatformSKU, error) {
	var m PlatformSKUModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND store_id = ? AND platform_code = ?", tenantID, storeID, platformCode).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &domain.PlatformSKU{
		ID: m.ID, TenantID: m.TenantID, SKUID: m.SKUID, StoreID: m.StoreID,
		PlatformCode: m.PlatformCode, PlatformSKU: m.PlatformSKUID, ASIN: m.ASIN, FNSKU: m.FNSKU,
		PlatformStatus: m.PlatformStatus, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}
