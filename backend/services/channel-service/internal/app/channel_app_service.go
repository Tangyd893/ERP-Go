package app

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/channel-service/internal/infra/repository"
)

// ChannelAppService 渠道应用服务
type ChannelAppService struct {
	repo        *repository.ChannelRepository
	orderClient *OrderServiceClient
	adapters    map[string]PlatformAdapter
}

func NewChannelAppService(repo *repository.ChannelRepository) *ChannelAppService {
	return &ChannelAppService{
		repo:     repo,
		adapters: make(map[string]PlatformAdapter),
	}
}

func (s *ChannelAppService) WithOrderClient(client *OrderServiceClient) *ChannelAppService {
	s.orderClient = client
	return s
}

func (s *ChannelAppService) RegisterAdapter(adapter PlatformAdapter) {
	s.adapters[adapter.PlatformCode()] = adapter
}

// ImportCSVOrders 导入 CSV 订单：解析 → 创建 OrderImportTask → 逐条调 order-service
func (s *ChannelAppService) ImportCSVOrders(ctx context.Context, tenantID, storeID, fileName string, csvData io.Reader) (*domain.OrderImportTask, int, int, error) {
	orders, err := ParseCSVOrders(csvData)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("CSV 解析失败: %w", err)
	}

	task := &domain.OrderImportTask{
		ID:             fmt.Sprintf("IM%d", time.Now().UnixNano()),
		TenantID:       tenantID,
		StoreID:        storeID,
		ImportType:     "csv",
		FileName:       fileName,
		IdempotencyKey: fmt.Sprintf("%s-%d", storeID, time.Now().Unix()),
		Status:         "processing",
		TotalRows:      len(orders),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreateImportTask(ctx, task); err != nil {
		return nil, 0, 0, fmt.Errorf("创建导入任务失败: %w", err)
	}

	success, failed := 0, 0
	if s.orderClient != nil {
		for _, order := range orders {
			if _, err := s.orderClient.CreateOrder(ctx, tenantID, order); err != nil {
				failed++
			} else {
				success++
			}
		}
	}

	task.Status = "completed"
	task.SuccessRows = success
	task.FailedRows = failed
	task.UpdatedAt = time.Now()

	return task, success, failed, nil
}

// FetchPlatformOrders 通过平台适配器拉取订单
func (s *ChannelAppService) FetchPlatformOrders(ctx context.Context, tenantID, storeID, platform string) ([]ImportedOrder, error) {
	adapter, ok := s.adapters[platform]
	if !ok {
		return nil, fmt.Errorf("未注册的平台适配器: %s", platform)
	}
	orders, err := adapter.FetchOrders(ctx, storeID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("拉取平台订单失败: %w", err)
	}
	return orders, nil
}

func (s *ChannelAppService) CreateStore(ctx context.Context, store *domain.Store) error {
	return s.repo.CreateStore(ctx, store)
}

func (s *ChannelAppService) ListStores(ctx context.Context, tenantID string) ([]*domain.Store, error) {
	return s.repo.ListStores(ctx, tenantID)
}

func (s *ChannelAppService) CreateImportTask(ctx context.Context, task *domain.OrderImportTask) error {
	return s.repo.CreateImportTask(ctx, task)
}

func (s *ChannelAppService) GetImportTask(ctx context.Context, idempotencyKey string) (*domain.OrderImportTask, error) {
	return s.repo.FindImportTaskByKey(ctx, idempotencyKey)
}

// CreateTrackingUpload 创建发货回传同步任务
func (s *ChannelAppService) CreateTrackingUpload(ctx context.Context, tenantID, storeID, trackingNo, carrierCode, orderNo string) (*domain.SyncTask, error) {
	task := &domain.SyncTask{
		ID:         fmt.Sprintf("TU-%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		StoreID:    storeID,
		TaskType:   "tracking_upload",
		Status:     "pending",
		TotalCount: 1,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.CreateSyncTask(ctx, task); err != nil {
		return nil, fmt.Errorf("创建回传任务失败: %w", err)
	}
	return task, nil
}
