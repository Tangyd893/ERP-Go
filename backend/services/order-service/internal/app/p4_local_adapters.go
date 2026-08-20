package app

import (
	"context"
	"fmt"
)

// LocalOrderStatusUpdater 在本服务内更新订单状态
type LocalOrderStatusUpdater struct {
	repo orderRepo
}

func NewLocalOrderStatusUpdater(repo orderRepo) *LocalOrderStatusUpdater {
	return &LocalOrderStatusUpdater{repo: repo}
}

func (u *LocalOrderStatusUpdater) UpdateOrderStatus(ctx context.Context, orderID, status string, metadata map[string]interface{}) error {
	tenantID := ""
	if metadata != nil {
		if v, ok := metadata["tenant_id"].(string); ok {
			tenantID = v
		}
	}
	if tenantID == "" {
		return fmt.Errorf("缺少租户上下文")
	}
	return u.repo.UpdateStatus(ctx, tenantID, orderID, status)
}
