package domain

import "time"

// SPU 标准产品单元（聚合根）
type SPU struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	CategoryID  string    `json:"category_id"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Images      []string  `json:"images"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SKU 库存单位（聚合根）
type SKU struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	SPUID        string    `json:"spu_id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Barcode      string    `json:"barcode"`
	SpecDesc     string    `json:"spec_desc"`
	Weight       float64   `json:"weight"`
	WeightUnit   string    `json:"weight_unit"`
	Length       float64   `json:"length"`
	Width        float64   `json:"width"`
	Height       float64   `json:"height"`
	LengthUnit   string    `json:"length_unit"`
	PurchasePrice float64  `json:"purchase_price"`
	SalePrice    float64   `json:"sale_price"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// VariantOption 变体选项（颜色、尺寸等）
type VariantOption struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	SortOrder int `json:"sort_order"`
}

// DeclarationInfo 海关申报信息
type DeclarationInfo struct {
	CNName        string  `json:"cn_name"`
	ENName        string  `json:"en_name"`
	HSCode        string  `json:"hs_code"`
	Material      string  `json:"material"`
	Usage         string  `json:"usage"`
	UnitPrice     float64 `json:"unit_price"`
	CustomsWeight float64 `json:"customs_weight"`
}

// PlatformSKU 平台 SKU 映射
type PlatformSKU struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	SKUID         string    `json:"sku_id"`
	StoreID       string    `json:"store_id"`
	PlatformCode  string    `json:"platform_code"`
	PlatformSKU   string    `json:"platform_sku"`
	ASIN          string    `json:"asin"`
	FNSKU         string    `json:"fnsku"`
	PlatformStatus string   `json:"platform_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
