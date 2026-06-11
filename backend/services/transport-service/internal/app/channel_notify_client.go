package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ChannelNotifyClient 渠道服务回传通知客户端
type ChannelNotifyClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewChannelNotifyClient 创建渠道通知客户端
func NewChannelNotifyClient(baseURL string) *ChannelNotifyClient {
	return &ChannelNotifyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// TrackingUploadRequest 发货回传请求
type TrackingUploadRequest struct {
	StoreID     string `json:"store_id"`
	TrackingNo  string `json:"tracking_no"`
	CarrierCode string `json:"carrier_code"`
	OrderNo     string `json:"order_no"`
}

// NotifyTrackingUpload 通知渠道服务创建 tracking_upload 同步任务
func (c *ChannelNotifyClient) NotifyTrackingUpload(ctx context.Context, tenantID string, req TrackingUploadRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化回传请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/channel/tracking-upload", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("创建回传请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-ID", tenantID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("通知渠道服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("渠道服务返回错误: HTTP %d", resp.StatusCode)
	}
	return nil
}
