package opensearchx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client 轻量 OpenSearch HTTP 客户端（T-604：Ping / EnsureIndex / IndexDoc）。
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewFromEnv 读取 OPENSEARCH_URL（默认 http://localhost:9200）；未配置则返回 nil。
func NewFromEnv() *Client {
	url := os.Getenv("OPENSEARCH_URL")
	if url == "" {
		url = os.Getenv("OPENSEARCH_HOST")
		if url == "" {
			return nil
		}
		if url[0] != 'h' {
			url = "http://" + url
		}
	}
	return &Client{
		baseURL:    url,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

// Ping 调用 GET / 检查集群可达。
func (c *Client) Ping(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("opensearch not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("opensearch status %d", resp.StatusCode)
	}
	return nil
}

// BaseURL 返回配置的集群地址。
func (c *Client) BaseURL() string {
	if c == nil {
		return ""
	}
	return c.baseURL
}
