package domain

import (
	"testing"
	"time"
)

// 创建测试用 SPU
func setupSPU() *SPU {
	return &SPU{
		ID:          "spu-001",
		TenantID:    "tenant-001",
		Name:        "测试商品",
		CategoryID:  "cat-001",
		Brand:       "测试品牌",
		Description: "这是一个测试商品",
		Images:      []string{"http://img.example.com/1.jpg", "http://img.example.com/2.jpg"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// 创建测试用 SKU
func setupSKU() *SKU {
	return &SKU{
		ID:            "sku-001",
		TenantID:      "tenant-001",
		SPUID:         "spu-001",
		Code:          "TSHIRT-RED-M",
		Name:          "红色T恤 M码",
		Barcode:       "6901234567890",
		SpecDesc:      "颜色:红色, 尺码:M",
		Weight:        0.2,
		WeightUnit:    "kg",
		Length:        30,
		Width:         20,
		Height:        2,
		LengthUnit:    "cm",
		PurchasePrice: 25.00,
		SalePrice:     99.00,
		Currency:      "CNY",
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// 创建测试用变体选项
func setupVariantOptions() []*VariantOption {
	return []*VariantOption{
		{ID: "vo-1", Name: "颜色", Value: "红色", SortOrder: 1},
		{ID: "vo-2", Name: "颜色", Value: "蓝色", SortOrder: 2},
		{ID: "vo-3", Name: "尺码", Value: "M", SortOrder: 1},
		{ID: "vo-4", Name: "尺码", Value: "L", SortOrder: 2},
	}
}

// 创建测试用平台 SKU
func setupPlatformSKU() *PlatformSKU {
	return &PlatformSKU{
		ID:             "psku-001",
		TenantID:       "tenant-001",
		SKUID:          "sku-001",
		StoreID:        "store-001",
		PlatformCode:   "amazon",
		PlatformSKU:    "AMZ-TSHIRT-RED-M",
		ASIN:           "B0123456789",
		FNSKU:          "X00123456789",
		PlatformStatus: "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// 创建测试用申报信息
func setupDeclarationInfo() *DeclarationInfo {
	return &DeclarationInfo{
		CNName:        "棉质红色T恤",
		ENName:        "Cotton Red T-Shirt",
		HSCode:        "6109100010",
		Material:      "纯棉",
		Usage:         "日常穿着",
		UnitPrice:     8.50,
		CustomsWeight: 0.25,
	}
}

// TestSPUCreation 测试 SPU 创建与基础字段
func TestSPUCreation(t *testing.T) {
	spu := setupSPU()

	if spu.ID == "" {
		t.Error("SPU ID不应为空")
	}
	if spu.Name == "" {
		t.Error("SPU 名称不应为空")
	}
	if spu.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if spu.Status != "active" {
		t.Errorf("默认状态应为 active，实际 %s", spu.Status)
	}
}

// TestSPUImages 测试 SPU 图片列表
func TestSPUImages(t *testing.T) {
	spu := setupSPU()

	if len(spu.Images) != 2 {
		t.Errorf("图片数量应为2，实际 %d", len(spu.Images))
	}

	// 追加图片
	spu.Images = append(spu.Images, "http://img.example.com/3.jpg")
	if len(spu.Images) != 3 {
		t.Errorf("追加后图片数量应为3，实际 %d", len(spu.Images))
	}

	// 空图片列表
	spu.Images = []string{}
	if len(spu.Images) != 0 {
		t.Errorf("清空后图片数量应为0，实际 %d", len(spu.Images))
	}
}

// TestSPUStatus 测试 SPU 状态值
func TestSPUStatus(t *testing.T) {
	validStatuses := []string{"active", "inactive", "discontinued", "draft"}

	spu := setupSPU()
	for _, s := range validStatuses {
		spu.Status = s
		if spu.Status != s {
			t.Errorf("状态应为 %s，实际 %s", s, spu.Status)
		}
	}
}

// TestSPUWithoutCategory 测试无分类的 SPU
func TestSPUWithoutCategory(t *testing.T) {
	spu := setupSPU()
	spu.CategoryID = ""

	if spu.CategoryID != "" {
		t.Error("未分类的 SPU 类别ID应为空")
	}
}

// TestSKUCreation 测试 SKU 创建与基础字段
func TestSKUCreation(t *testing.T) {
	sku := setupSKU()

	if sku.ID == "" {
		t.Error("SKU ID不应为空")
	}
	if sku.SPUID == "" {
		t.Error("SKU 的 SPUID不应为空")
	}
	if sku.Code == "" {
		t.Error("SKU 编码不应为空")
	}
	if sku.Barcode == "" {
		t.Error("条形码不应为空")
	}
}

// TestSKUDimensions 测试 SKU 尺寸与重量
func TestSKUDimensions(t *testing.T) {
	sku := setupSKU()

	if sku.Weight <= 0 {
		t.Error("重量应大于0")
	}
	if sku.Length <= 0 || sku.Width <= 0 || sku.Height <= 0 {
		t.Error("长宽高均应大于0")
	}
	if sku.WeightUnit != "kg" {
		t.Errorf("重量单位应为 kg，实际 %s", sku.WeightUnit)
	}
	if sku.LengthUnit != "cm" {
		t.Errorf("长度单位应为 cm，实际 %s", sku.LengthUnit)
	}
}

// TestSKUPricing 测试 SKU 价格字段
func TestSKUPricing(t *testing.T) {
	sku := setupSKU()

	if sku.PurchasePrice <= 0 {
		t.Error("采购价应大于0")
	}
	if sku.SalePrice <= 0 {
		t.Error("销售价应大于0")
	}
	if sku.SalePrice <= sku.PurchasePrice {
		t.Error("销售价应大于采购价才有毛利")
	}
}

// TestSKUZeroDimension 测试 SKU 零尺寸边界
func TestSKUZeroDimension(t *testing.T) {
	sku := &SKU{
		ID:     "sku-zero-size",
		Weight: 0,
		Length: 0,
		Width:  0,
		Height: 0,
	}

	// 零尺寸在领域层不应报错，由应用层处理校验
	if sku.Weight != 0 || sku.Length != 0 || sku.Width != 0 || sku.Height != 0 {
		t.Error("零尺寸 SKU 维度值应为0")
	}
}

// TestSKUStatus 测试 SKU 状态值
func TestSKUStatus(t *testing.T) {
	validStatuses := []string{"active", "inactive", "discontinued"}

	sku := setupSKU()
	for _, s := range validStatuses {
		sku.Status = s
		if sku.Status != s {
			t.Errorf("状态应为 %s，实际 %s", s, sku.Status)
		}
	}
}

// TestSKUCurrency 测试货币单位
func TestSKUCurrency(t *testing.T) {
	currencies := []string{"CNY", "USD", "EUR", "JPY"}

	sku := setupSKU()
	for _, c := range currencies {
		sku.Currency = c
		if sku.Currency != c {
			t.Errorf("货币应为 %s，实际 %s", c, sku.Currency)
		}
	}
}

// TestVariantOptionSortOrder 测试变体选项排序
func TestVariantOptionSortOrder(t *testing.T) {
	options := setupVariantOptions()

	if len(options) != 4 {
		t.Errorf("变体选项数量应为4，实际 %d", len(options))
	}

	// 验证 SortOrder 有效
	for i, opt := range options {
		if opt.SortOrder <= 0 {
			t.Errorf("选项 %d 排序值应大于0，实际 %d", i, opt.SortOrder)
		}
		if opt.ID == "" {
			t.Errorf("选项 %d ID不应为空", i)
		}
		if opt.Name == "" {
			t.Errorf("选项 %d 名称不应为空", i)
		}
	}
}

// TestVariantOptionFieldValues 测试变体选项字段赋值
func TestVariantOptionFieldValues(t *testing.T) {
	opt := &VariantOption{
		ID:        "vo-test",
		Name:      "尺寸",
		Value:     "XL",
		SortOrder: 5,
	}

	if opt.Name != "尺寸" {
		t.Errorf("选项名应为 尺寸，实际 %s", opt.Name)
	}
	if opt.Value != "XL" {
		t.Errorf("选项值应为 XL，实际 %s", opt.Value)
	}
	if opt.SortOrder != 5 {
		t.Errorf("排序应为5，实际 %d", opt.SortOrder)
	}
}

// TestDeclarationInfoCreation 测试海关申报信息创建
func TestDeclarationInfoCreation(t *testing.T) {
	info := setupDeclarationInfo()

	if info.CNName == "" {
		t.Error("中文名称不应为空")
	}
	if info.ENName == "" {
		t.Error("英文名称不应为空")
	}
	if info.HSCode == "" {
		t.Error("HS编码不应为空")
	}
	if info.Material == "" {
		t.Error("材质不应为空")
	}
	if info.Usage == "" {
		t.Error("用途不应为空")
	}
}

// TestDeclarationInfoPricing 测试申报信息价格
func TestDeclarationInfoPricing(t *testing.T) {
	info := setupDeclarationInfo()

	if info.UnitPrice <= 0 {
		t.Error("申报单价应大于0")
	}
	if info.CustomsWeight <= 0 {
		t.Error("报关重量应大于0")
	}
}

// TestDeclarationInfoHSCodeLength 测试 HS 编码格式
func TestDeclarationInfoHSCodeLength(t *testing.T) {
	tests := []struct {
		name   string
		hsCode string
	}{
		{"10位HS编码", "6109100010"},
		{"8位HS编码", "61091000"},
		{"6位HS编码", "610910"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &DeclarationInfo{
				CNName:   "测试商品",
				ENName:   "Test Product",
				HSCode:   tt.hsCode,
				Material: "测试材质",
				Usage:    "测试用途",
			}
			if info.HSCode != tt.hsCode {
				t.Errorf("HS编码应为 %s，实际 %s", tt.hsCode, info.HSCode)
			}
		})
	}
}

// TestPlatformSKUCreation 测试平台 SKU 映射创建
func TestPlatformSKUCreation(t *testing.T) {
	psku := setupPlatformSKU()

	if psku.ID == "" {
		t.Error("平台SKU ID不应为空")
	}
	if psku.SKUID == "" {
		t.Error("关联的SKU ID不应为空")
	}
	if psku.StoreID == "" {
		t.Error("店铺ID不应为空")
	}
	if psku.PlatformCode == "" {
		t.Error("平台编码不应为空")
	}
	if psku.PlatformSKU == "" {
		t.Error("平台SKU编码不应为空")
	}
}

// TestPlatformSKUAsinFNSKU 测试亚马逊特有字段
func TestPlatformSKUAsinFNSKU(t *testing.T) {
	psku := setupPlatformSKU()

	if psku.ASIN == "" {
		t.Error("亚马逊ASIN不应为空")
	}
	if psku.FNSKU == "" {
		t.Error("FNSKU不应为空")
	}
}

// TestPlatformSKUPlatformCode 测试不同平台编码映射
func TestPlatformSKUPlatformCode(t *testing.T) {
	platforms := []string{"amazon", "shopee", "lazada", "tiktok", "ebay"}

	for _, p := range platforms {
		t.Run("平台_"+p, func(t *testing.T) {
			psku := &PlatformSKU{
				ID:             "psku-" + p,
				TenantID:       "tenant-001",
				SKUID:          "sku-001",
				StoreID:        "store-001",
				PlatformCode:   p,
				PlatformSKU:    p + "-SKU-001",
				PlatformStatus: "active",
			}
			if psku.PlatformCode != p {
				t.Errorf("平台编码应为 %s，实际 %s", p, psku.PlatformCode)
			}
		})
	}
}

// TestPlatformSKUStatus 测试平台 SKU 状态
func TestPlatformSKUStatus(t *testing.T) {
	statuses := []string{"active", "inactive", "out_of_stock", "discontinued"}

	psku := setupPlatformSKU()
	for _, s := range statuses {
		psku.PlatformStatus = s
		if psku.PlatformStatus != s {
			t.Errorf("平台状态应为 %s，实际 %s", s, psku.PlatformStatus)
		}
	}
}

// TestPlatformSKUWithoutAsin 测试非亚马逊平台无需 ASIN
func TestPlatformSKUWithoutAsin(t *testing.T) {
	psku := &PlatformSKU{
		ID:             "psku-shopee-001",
		TenantID:       "tenant-001",
		SKUID:          "sku-001",
		StoreID:        "store-001",
		PlatformCode:   "shopee",
		PlatformSKU:    "SHOPEE-SKU-001",
		PlatformStatus: "active",
	}

	// 非亚马逊平台 ASIN 和 FNSKU 可为空
	if psku.ASIN != "" {
		t.Error("非亚马逊平台 ASIN 应为空")
	}
	if psku.FNSKU != "" {
		t.Error("非亚马逊平台 FNSKU 应为空")
	}
}

// 测试重量单位变更
func TestSKUWeightUnitConversionAwareness(t *testing.T) {
	tests := []struct {
		name       string
		weightUnit string
	}{
		{"千克", "kg"},
		{"克", "g"},
		{"磅", "lb"},
		{"盎司", "oz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sku := setupSKU()
			sku.WeightUnit = tt.weightUnit
			if sku.WeightUnit != tt.weightUnit {
				t.Errorf("重量单位应为 %s，实际 %s", tt.weightUnit, sku.WeightUnit)
			}
		})
	}
}
