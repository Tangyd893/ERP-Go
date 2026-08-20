package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/order-service/internal/domain"
	opensearchx "github.com/Tangyd893/ERP-Go/backend/shared/opensearch"
)

// OrderSearchIndexer 将订单写入 OpenSearch（T-604 业务挂接试点）。
type OrderSearchIndexer struct {
	client *opensearchx.Client
}

func NewOrderSearchIndexer(c *opensearchx.Client) *OrderSearchIndexer {
	if c == nil {
		return nil
	}
	return &OrderSearchIndexer{client: c}
}

// IndexOrder best-effort 覆盖写入；失败由调用方忽略，不影响主流程。
func (i *OrderSearchIndexer) IndexOrder(ctx context.Context, order *domain.SalesOrder) error {
	if i == nil || i.client == nil || order == nil {
		return nil
	}
	doc := opensearchx.OrderDoc{
		OrderID:         order.ID,
		TenantID:        order.TenantID,
		PlatformOrderNo: order.PlatformOrderNo,
		Status:          string(order.Status),
		TotalAmount:     order.TotalAmount,
	}
	return i.client.IndexDoc(ctx, opensearchx.IndexOrders, order.ID, doc)
}
