package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/transport-service/internal/infra/repository"
)

// LabelAdapter 面单适配器接口（首期 Mock，后续对接真实物流商 API）
type LabelAdapter interface {
	CreateLabel(ctx context.Context, req LabelRequest) (*LabelResponse, error)
}

// LabelRequest 面单创建请求
type LabelRequest struct {
	ShipmentID  string
	OrderNo     string
	CarrierCode string
	ServiceCode string
	Weight      float64
	ToAddress   AddressInfo
}

// AddressInfo 地址信息
type AddressInfo struct {
	Name       string
	Phone      string
	Country    string
	State      string
	City       string
	District   string
	StreetLine string
	PostalCode string
}

// LabelResponse 面单响应
type LabelResponse struct {
	TrackingNo string
	LabelURL   string
	CarrierCode string
	ServiceCode string
}

// MockLabelAdapter Mock 面单适配器（开发用）
type MockLabelAdapter struct{}

func (m *MockLabelAdapter) CreateLabel(ctx context.Context, req LabelRequest) (*LabelResponse, error) {
	trackingNo := fmt.Sprintf("MOCK%d", time.Now().UnixNano()%100000000000)
	return &LabelResponse{
		TrackingNo:  trackingNo,
		LabelURL:    fmt.Sprintf("https://mock-carrier.example.com/labels/%s.pdf", trackingNo),
		CarrierCode: req.CarrierCode,
		ServiceCode: req.ServiceCode,
	}, nil
}

// ── TransportAppService ────────────────────────────────

type TransportAppService struct {
	repo          *repository.TransportRepository
	labelAdapter  LabelAdapter
	channelClient *ChannelNotifyClient
}

func NewTransportAppService(repo *repository.TransportRepository) *TransportAppService {
	return &TransportAppService{repo: repo, labelAdapter: &MockLabelAdapter{}}
}

// WithLabelAdapter 注入自定义面单适配器
func (s *TransportAppService) WithLabelAdapter(adapter LabelAdapter) *TransportAppService {
	s.labelAdapter = adapter
	return s
}

// WithChannelClient 注入渠道通知客户端
func (s *TransportAppService) WithChannelClient(client *ChannelNotifyClient) *TransportAppService {
	s.channelClient = client
	return s
}

func (s *TransportAppService) ListCarriers(ctx context.Context, tenantID string) ([]*domain.Carrier, error) {
	return s.repo.ListCarriers(ctx, tenantID)
}

// MatchCarrier 物流渠道匹配：按重量+目的地匹配最优物流产品
func (s *TransportAppService) MatchCarrier(ctx context.Context, tenantID string, weight float64, country string) (*MatchResult, error) {
	rules, err := s.repo.ListShippingRules(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询物流规则失败: %w", err)
	}
	for _, rule := range rules {
		// 重量范围匹配（0 表示不限）
		if rule.MinWeight > 0 && weight < rule.MinWeight {
			continue
		}
		if rule.MaxWeight > 0 && weight > rule.MaxWeight {
			continue
		}
		// 国家匹配（空列表表示不限）
		if len(rule.CountryCodes) > 0 {
			matched := false
			for _, c := range rule.CountryCodes {
				if strings.EqualFold(c, country) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		// 命中最高优先级规则
		svc, err := s.repo.FindCarrierService(ctx, rule.CarrierServiceID)
		if err != nil {
			continue
		}
		return &MatchResult{
			RuleID:        rule.ID,
			RuleName:      rule.Name,
			CarrierService: svc,
		}, nil
	}
	return nil, fmt.Errorf("未找到匹配的物流渠道（重量=%.0fg, 国家=%s）", weight, country)
}

// MatchResult 匹配结果
type MatchResult struct {
	RuleID         string
	RuleName       string
	CarrierService *domain.CarrierService
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

// GenerateLabel 生成面单：调用适配器 + 更新发运状态
func (s *TransportAppService) GenerateLabel(ctx context.Context, shipmentID, orderNo, toCountry string, weight float64, addr AddressInfo) (*LabelResponse, error) {
	shipment, err := s.repo.FindShipment(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("发运单不存在: %w", err)
	}
	if s.labelAdapter == nil {
		return nil, fmt.Errorf("面单适配器未配置")
	}
	resp, err := s.labelAdapter.CreateLabel(ctx, LabelRequest{
		ShipmentID:  shipmentID,
		OrderNo:     orderNo,
		CarrierCode: shipment.CarrierCode,
		ServiceCode: shipment.ServiceCode,
		Weight:      weight,
		ToAddress:   addr,
	})
	if err != nil {
		return nil, fmt.Errorf("生成面单失败: %w", err)
	}
	if err := s.repo.UpdateShipmentStatus(ctx, shipmentID, string(domain.ShipmentLabeled), resp.TrackingNo); err != nil {
		return nil, err
	}

	// 通知渠道服务创建发货回传任务
	if s.channelClient != nil {
		_ = s.channelClient.NotifyTrackingUpload(ctx, shipment.TenantID, TrackingUploadRequest{
			StoreID:     shipment.OrderID, // 通过 OrderID 关联到店铺
			TrackingNo:  resp.TrackingNo,
			CarrierCode: resp.CarrierCode,
			OrderNo:     orderNo,
		})
	}

	return resp, nil
}

func (s *TransportAppService) CreateLabel(ctx context.Context, id, trackingNo, labelURL string) error {
	return s.repo.UpdateShipmentStatus(ctx, id, string(domain.ShipmentLabeled), trackingNo)
}
