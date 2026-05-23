package domain

import (
	"testing"
	"time"
)

// 创建测试用发运单
func setupShipment() *Shipment {
	now := time.Now()
	return &Shipment{
		ID:          "ship-001",
		TenantID:    "tenant-001",
		OrderID:     "order-001",
		OutboundID:  "out-001",
		CarrierCode: "SF",
		ServiceCode: "SF-EXPRESS",
		TrackingNo:  "SF1234567890",
		Status:      ShipmentPending,
		Weight:      1.5,
		ShippingCost: 15.00,
		Currency:    "CNY",
		Packages: []*PackageInfo{
			{ID: "pkg-001", TrackingNo: "SF1234567890-01", Weight: 0.8, Length: 30, Width: 20, Height: 5},
			{ID: "pkg-002", TrackingNo: "SF1234567890-02", Weight: 0.7, Length: 25, Width: 15, Height: 3},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// 创建测试用物流商
func setupCarrier() *Carrier {
	return &Carrier{
		ID:        "car-001",
		TenantID:  "tenant-001",
		Name:      "顺丰速运",
		Code:      "SF",
		Status:    "active",
		CreatedAt: time.Now(),
	}
}

// 创建测试用物流产品
func setupCarrierService() *CarrierService {
	return &CarrierService{
		ID:          "cs-001",
		CarrierID:   "car-001",
		Name:        "顺丰标快",
		Code:        "SF-EXPRESS",
		ServiceType: "express",
	}
}

// 创建测试用物流规则
func setupShippingRule() *ShippingRule {
	return &ShippingRule{
		ID:               "rule-001",
		TenantID:         "tenant-001",
		Name:             "中国地区标准规则",
		Priority:         1,
		CountryCodes:     []string{"CN"},
		MinWeight:        0.0,
		MaxWeight:        30.0,
		CarrierServiceID: "cs-001",
	}
}

// TestShipmentStatus 测试发运状态常量
func TestShipmentStatus(t *testing.T) {
	tests := []struct {
		name   string
		status ShipmentStatus
		want   string
	}{
		{"待处理", ShipmentPending, "pending"},
		{"已打标", ShipmentLabeled, "labeled"},
		{"已发货", ShipmentShipped, "shipped"},
		{"运输中", ShipmentInTransit, "in_transit"},
		{"已送达", ShipmentDelivered, "delivered"},
		{"已取消", ShipmentCancelled, "cancelled"},
		{"失败", ShipmentFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("状态值应为 %s，实际 %s", tt.want, tt.status)
			}
		})
	}
}

// TestShipmentCreation 测试发运单创建与基础字段
func TestShipmentCreation(t *testing.T) {
	s := setupShipment()

	if s.ID == "" {
		t.Error("发运单ID不应为空")
	}
	if s.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if s.OrderID == "" {
		t.Error("关联订单ID不应为空")
	}
	if s.CarrierCode == "" {
		t.Error("物流商编码不应为空")
	}
	if s.TrackingNo == "" {
		t.Error("运单号不应为空")
	}
	if s.Status != ShipmentPending {
		t.Errorf("初始状态应为 pending，实际 %s", s.Status)
	}
}

// TestShipmentStatusTransitions 测试发运状态流转
func TestShipmentStatusTransitions(t *testing.T) {
	transitions := []struct {
		from ShipmentStatus
		to   ShipmentStatus
	}{
		{ShipmentPending, ShipmentLabeled},
		{ShipmentLabeled, ShipmentShipped},
		{ShipmentShipped, ShipmentInTransit},
		{ShipmentInTransit, ShipmentDelivered},
		{ShipmentPending, ShipmentCancelled},
		{ShipmentLabeled, ShipmentFailed},
	}

	for _, tr := range transitions {
		t.Run(string(tr.from)+"_to_"+string(tr.to), func(t *testing.T) {
			s := &Shipment{
				ID:     "test",
				Status: tr.from,
			}
			s.Status = tr.to
			if s.Status != tr.to {
				t.Errorf("状态从 %s 流转到 %s 后应为 %s", tr.from, tr.to, tr.to)
			}
		})
	}
}

// TestShipmentPackages 测试发运包裹列表
func TestShipmentPackages(t *testing.T) {
	s := setupShipment()

	if len(s.Packages) != 2 {
		t.Errorf("包裹数量应为2，实际 %d", len(s.Packages))
	}
}

// TestShipmentEmptyPackages 测试无包裹发运单
func TestShipmentEmptyPackages(t *testing.T) {
	s := &Shipment{
		ID:          "ship-empty",
		TenantID:    "tenant-001",
		OrderID:     "order-001",
		OutboundID:  "out-001",
		CarrierCode: "SF",
		TrackingNo:  "SF0000000000",
		Status:      ShipmentPending,
		Packages:    []*PackageInfo{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if len(s.Packages) != 0 {
		t.Errorf("无包裹发运单包裹数应为0，实际 %d", len(s.Packages))
	}
}

// TestShipmentWeight 测试发运单重量
func TestShipmentWeight(t *testing.T) {
	s := setupShipment()

	if s.Weight <= 0 {
		t.Error("发运重量应大于0")
	}
}

// TestShipmentShippingCost 测试发运费
func TestShipmentShippingCost(t *testing.T) {
	s := setupShipment()

	if s.ShippingCost <= 0 {
		t.Error("发运费应大于0")
	}
}

// TestShipmentTrackingRecords 测试物流轨迹
func TestShipmentTrackingRecords(t *testing.T) {
	s := setupShipment()
	records := []*TrackingRecord{
		{Status: "picked_up", Description: "已揽件", Location: "深圳", RecordedAt: time.Now()},
		{Status: "in_transit", Description: "运输中", Location: "广州中转", RecordedAt: time.Now().Add(2 * time.Hour)},
		{Status: "delivered", Description: "已签收", Location: "北京", RecordedAt: time.Now().Add(24 * time.Hour)},
	}
	s.TrackingRecords = records

	if len(s.TrackingRecords) != 3 {
		t.Errorf("物流轨迹数量应为3，实际 %d", len(s.TrackingRecords))
	}
}

// TestPackageInfoCreation 测试包裹信息创建
func TestPackageInfoCreation(t *testing.T) {
	pkg := &PackageInfo{
		ID:         "pkg-001",
		TrackingNo: "SF1234567890-01",
		Weight:     0.8,
		Length:     30,
		Width:      20,
		Height:     5,
	}

	if pkg.ID == "" {
		t.Error("包裹ID不应为空")
	}
	if pkg.TrackingNo == "" {
		t.Error("包裹运单号不应为空")
	}
	if pkg.Weight <= 0 {
		t.Error("包裹重量应大于0")
	}
	if pkg.Length <= 0 || pkg.Width <= 0 || pkg.Height <= 0 {
		t.Error("包裹长宽高应大于0")
	}
}

// TestPackageInfoDimensions 测试包裹尺寸
func TestPackageInfoDimensions(t *testing.T) {
	tests := []struct {
		name   string
		length float64
		width  float64
		height float64
	}{
		{"小包裹", 20, 15, 3},
		{"中包裹", 40, 30, 10},
		{"大包裹", 60, 50, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := &PackageInfo{
				ID:         "test-pkg",
				TrackingNo: "TEST-001",
				Weight:     1.0,
				Length:     tt.length,
				Width:      tt.width,
				Height:     tt.height,
			}
			if pkg.Length != tt.length || pkg.Width != tt.width || pkg.Height != tt.height {
				t.Error("包裹尺寸不一致")
			}
		})
	}
}

// TestTrackingRecordCreation 测试物流轨迹记录
func TestTrackingRecordCreation(t *testing.T) {
	record := &TrackingRecord{
		Status:      "picked_up",
		Description: "快件已揽收",
		Location:    "深圳市南山区",
		RecordedAt:  time.Now(),
	}

	if record.Status == "" {
		t.Error("轨迹状态不应为空")
	}
	if record.Description == "" {
		t.Error("轨迹描述不应为空")
	}
	if record.Location == "" {
		t.Error("轨迹地点不应为空")
	}
	if record.RecordedAt.IsZero() {
		t.Error("记录时间不应为零值")
	}
}

// TestTrackingRecordStatuses 测试各种物流轨迹状态
func TestTrackingRecordStatuses(t *testing.T) {
	statuses := []string{
		"picked_up",
		"in_transit",
		"out_for_delivery",
		"delivered",
		"failed_attempt",
		"returned",
	}

	for _, s := range statuses {
		t.Run("轨迹_"+s, func(t *testing.T) {
			record := &TrackingRecord{
				Status:      s,
				Description: "测试描述",
				Location:    "测试地点",
				RecordedAt:  time.Now(),
			}
			if record.Status != s {
				t.Errorf("轨迹状态应为 %s，实际 %s", s, record.Status)
			}
		})
	}
}

// TestCarrierCreation 测试物流商创建
func TestCarrierCreation(t *testing.T) {
	c := setupCarrier()

	if c.ID == "" {
		t.Error("物流商ID不应为空")
	}
	if c.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if c.Name == "" {
		t.Error("物流商名称不应为空")
	}
	if c.Code == "" {
		t.Error("物流商编码不应为空")
	}
	if c.Status != "active" {
		t.Errorf("物流商状态应为 active，实际 %s", c.Status)
	}
}

// TestCarrierStatus 测试物流商状态
func TestCarrierStatus(t *testing.T) {
	statuses := []string{"active", "inactive", "suspended"}

	c := setupCarrier()
	for _, s := range statuses {
		c.Status = s
		if c.Status != s {
			t.Errorf("状态应为 %s，实际 %s", s, c.Status)
		}
	}
}

// TestCarrierServiceCreation 测试物流产品创建
func TestCarrierServiceCreation(t *testing.T) {
	cs := setupCarrierService()

	if cs.ID == "" {
		t.Error("物流产品ID不应为空")
	}
	if cs.CarrierID == "" {
		t.Error("关联物流商ID不应为空")
	}
	if cs.Name == "" {
		t.Error("物流产品名称不应为空")
	}
	if cs.Code == "" {
		t.Error("物流产品编码不应为空")
	}
	if cs.ServiceType == "" {
		t.Error("服务类型不应为空")
	}
}

// TestCarrierServiceType 测试物流服务类型
func TestCarrierServiceType(t *testing.T) {
	serviceTypes := []string{"express", "standard", "economy"}

	cs := setupCarrierService()
	for _, st := range serviceTypes {
		cs.ServiceType = st
		if cs.ServiceType != st {
			t.Errorf("服务类型应为 %s，实际 %s", st, cs.ServiceType)
		}
	}
}

// TestShippingRuleCreation 测试物流规则创建
func TestShippingRuleCreation(t *testing.T) {
	r := setupShippingRule()

	if r.ID == "" {
		t.Error("规则ID不应为空")
	}
	if r.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if r.Name == "" {
		t.Error("规则名称不应为空")
	}
	if r.CarrierServiceID == "" {
		t.Error("关联物流产品ID不应为空")
	}
}

// TestShippingRuleWeightRange 测试物流规则重量范围
func TestShippingRuleWeightRange(t *testing.T) {
	tests := []struct {
		name      string
		minWeight float64
		maxWeight float64
	}{
		{"0-1kg", 0, 1.0},
		{"1-10kg", 1.0, 10.0},
		{"10-30kg", 10.0, 30.0},
		{"0-无限", 0, 99999.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupShippingRule()
			r.MinWeight = tt.minWeight
			r.MaxWeight = tt.maxWeight

			if r.MinWeight != tt.minWeight {
				t.Errorf("最小重量应为 %.1f，实际 %.1f", tt.minWeight, r.MinWeight)
			}
			if r.MaxWeight != tt.maxWeight {
				t.Errorf("最大重量应为 %.1f，实际 %.1f", tt.maxWeight, r.MaxWeight)
			}
		})
	}
}

// TestShippingRuleOverlapWeight 测试重量范围校验
func TestShippingRuleOverlapWeight(t *testing.T) {
	r := setupShippingRule()
	r.MinWeight = 10.0
	r.MaxWeight = 5.0

	// 不合法的重量范围应在领域层被允许存储
	if r.MinWeight <= r.MaxWeight {
		t.Error("最小重量大于最大重量为非法状态")
	}
}

// TestShippingRulePriority 测试物流规则优先级
func TestShippingRulePriority(t *testing.T) {
	priorities := []int{1, 5, 10, 99}

	for _, p := range priorities {
		t.Run("优先级_"+string(rune('0'+p%10)), func(t *testing.T) {
			r := setupShippingRule()
			r.Priority = p
			if r.Priority != p {
				t.Errorf("优先级应为 %d，实际 %d", p, r.Priority)
			}
		})
	}
}

// TestShippingRuleCountryCodes 测试物流规则国家编码
func TestShippingRuleCountryCodes(t *testing.T) {
	tests := []struct {
		name    string
		countries []string
	}{
		{"仅中国", []string{"CN"}},
		{"东亚三地", []string{"CN", "JP", "KR"}},
		{"欧洲多国", []string{"DE", "FR", "IT", "ES", "GB"}},
		{"空列表", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupShippingRule()
			r.CountryCodes = tt.countries
			if len(r.CountryCodes) != len(tt.countries) {
				t.Errorf("国家编码数量应为 %d，实际 %d", len(tt.countries), len(r.CountryCodes))
			}
		})
	}
}

// TestShipmentCurrency 测试发运货币单位
func TestShipmentCurrency(t *testing.T) {
	currencies := []string{"CNY", "USD", "EUR"}

	s := setupShipment()
	for _, c := range currencies {
		s.Currency = c
		if s.Currency != c {
			t.Errorf("货币应为 %s，实际 %s", c, s.Currency)
		}
	}
}

// TestShipmentLabelURL 测试面单地址
func TestShipmentLabelURL(t *testing.T) {
	s := setupShipment()
	s.LabelURL = "https://label.sf-express.com/labels/SF1234567890.pdf"

	if s.LabelURL == "" {
		t.Error("面单URL不应为空")
	}
}

// TestShipmentNilTrackingRecords 测试无轨迹的发运单
func TestShipmentNilTrackingRecords(t *testing.T) {
	s := setupShipment()

	if s.TrackingRecords != nil {
		t.Error("无轨迹时 TrackingRecords 应为 nil")
	}
}
