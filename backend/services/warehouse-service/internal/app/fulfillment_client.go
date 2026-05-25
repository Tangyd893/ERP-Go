package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/shared/workflows"
)

// OrderFulfillmentClient 通知 Order 服务处理出库完成
type OrderFulfillmentClient struct {
	orderURL string
	client   *http.Client
}

func NewOrderFulfillmentClient(orderURL string) *OrderFulfillmentClient {
	return &OrderFulfillmentClient{
		orderURL: orderURL,
		client:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *OrderFulfillmentClient) NotifyOutboundShipped(ctx context.Context, data workflows.OutboundShippedData) error {
	if c.orderURL == "" {
		return fmt.Errorf("未配置 ORDER_SERVICE_URL")
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.orderURL+"/api/v1/order/fulfillment/outbound-shipped", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("通知订单服务失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("订单履约回调失败: HTTP %d", resp.StatusCode)
	}
	return nil
}
