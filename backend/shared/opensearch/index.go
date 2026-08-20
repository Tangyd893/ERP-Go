package opensearchx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	IndexProducts = "erp-products"
	IndexOrders   = "erp-orders"
)

// IndexDoc 写入/覆盖文档（PUT /{index}/_doc/{id}）。
func (c *Client) IndexDoc(ctx context.Context, index, id string, doc any) error {
	if c == nil {
		return fmt.Errorf("opensearch not configured")
	}
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/_doc/%s", strings.TrimRight(c.baseURL, "/"), index, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch index status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// EnsureIndex 若不存在则创建空 mapping 索引（idempotent）。
func (c *Client) EnsureIndex(ctx context.Context, index string) error {
	if c == nil {
		return fmt.Errorf("opensearch not configured")
	}
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), index)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(`{}`))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create index status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// ProductDoc / OrderDoc 索引文档草案字段。
type ProductDoc struct {
	SKUID    string `json:"sku_id"`
	SKUCode  string `json:"sku_code"`
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
}

type OrderDoc struct {
	OrderID         string  `json:"order_id"`
	TenantID        string  `json:"tenant_id"`
	PlatformOrderNo string  `json:"platform_order_no"`
	Status          string  `json:"status"`
	TotalAmount     float64 `json:"total_amount"`
}
