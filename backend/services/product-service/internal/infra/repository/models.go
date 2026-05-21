package repository

import "time"

type SPUModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	TenantID    string    `gorm:"column:tenant_id;index"`
	Name        string    `gorm:"column:name"`
	CategoryID  string    `gorm:"column:category_id"`
	Brand       string    `gorm:"column:brand"`
	Description string    `gorm:"column:description"`
	MainImage   string    `gorm:"column:main_image"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (SPUModel) TableName() string { return "spus" }

type SKUModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	TenantID      string    `gorm:"column:tenant_id;index"`
	SPUID         string    `gorm:"column:spu_id;index"`
	SKUCode       string    `gorm:"column:sku_code;unique"`
	Barcode       string    `gorm:"column:barcode"`
	SpecDesc      string    `gorm:"column:spec_desc"`
	WeightGram    int       `gorm:"column:weight_gram"`
	LengthCm      float64   `gorm:"column:length_cm"`
	WidthCm       float64   `gorm:"column:width_cm"`
	HeightCm      float64   `gorm:"column:height_cm"`
	PurchasePrice float64   `gorm:"column:purchase_price"`
	SellingPrice  float64   `gorm:"column:selling_price"`
	Currency      string    `gorm:"column:currency"`
	Status        string    `gorm:"column:status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (SKUModel) TableName() string { return "skus" }

type PlatformSKUModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	TenantID       string    `gorm:"column:tenant_id"`
	SKUID          string    `gorm:"column:sku_id;index"`
	StoreID        string    `gorm:"column:store_id;index"`
	PlatformCode   string    `gorm:"column:platform_code"`
	PlatformSKUID  string    `gorm:"column:platform_sku_id"`
	ASIN           string    `gorm:"column:asin"`
	FNSKU          string    `gorm:"column:fnsku"`
	PlatformStatus string    `gorm:"column:platform_status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (PlatformSKUModel) TableName() string { return "platform_skus" }
