// Package httpclient 提供内部服务间调用的共享 HTTP 客户端。
// 封装超时、认证头注入、统一响应解析。
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 内部服务 HTTP 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
	tenantID   string
}

// New 创建内部 HTTP 客户端
// baseURL 为服务地址，如 "http://localhost:8085"
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WithTenant 注入租户 ID（用于多租户数据隔离）
func (c *Client) WithTenant(tenantID string) *Client {
	c.tenantID = tenantID
	return c
}

// Resp 统一响应
type Resp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// PageResp 分页响应 data 结构
type PageResp struct {
	List       json.RawMessage `json:"list"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// GetPageTotal 调用 list 端点（page_size=1），仅提取 total 计数
func (c *Client) GetPageTotal(ctx context.Context, path string) (int64, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("httpclient: 创建请求失败: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("httpclient: 请求 %s 失败: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("httpclient: 读取响应失败: %w", err)
	}

	var r Resp
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, fmt.Errorf("httpclient: 解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return 0, fmt.Errorf("httpclient: 服务返回错误 code=%d msg=%s", r.Code, r.Message)
	}

	var page PageResp
	if err := json.Unmarshal(r.Data, &page); err != nil {
		return 0, fmt.Errorf("httpclient: 解析分页数据失败: %w", err)
	}
	return page.Total, nil
}

// GetList 调用 list 端点并返回原始 JSON list
func (c *Client) GetList(ctx context.Context, path string) (json.RawMessage, int64, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("httpclient: 创建请求失败: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("httpclient: 请求 %s 失败: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("httpclient: 读取响应失败: %w", err)
	}

	var r Resp
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, 0, fmt.Errorf("httpclient: 解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return nil, 0, fmt.Errorf("httpclient: 服务返回错误 code=%d msg=%s", r.Code, r.Message)
	}

	var page PageResp
	if err := json.Unmarshal(r.Data, &page); err != nil {
		return nil, 0, fmt.Errorf("httpclient: 解析分页数据失败: %w", err)
	}
	return page.List, page.Total, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if c.tenantID != "" {
		req.Header.Set("X-Tenant-ID", c.tenantID)
	}
}
