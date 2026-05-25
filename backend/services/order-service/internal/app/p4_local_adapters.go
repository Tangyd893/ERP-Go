package app

import (
	"context"
)

// LocalOrderStatusUpdater 在本服务内更新订单状态
type LocalOrderStatusUpdater struct {
	repo orderRepo
}

func NewLocalOrderStatusUpdater(repo orderRepo) *LocalOrderStatusUpdater {
	return &LocalOrderStatusUpdater{repo: repo}
}

func (u *LocalOrderStatusUpdater) UpdateOrderStatus(ctx context.Context, orderID, status string, metadata map[string]interface{}) error {
	_ = metadata
	return u.repo.UpdateStatus(ctx, orderID, status)
}
