package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"
)

// HTTPStockLockAdapter 通过 HTTP 调用 Inventory 服务锁定库存
type HTTPStockLockAdapter struct {
	inventoryURL string
	client       *http.Client
}

func NewHTTPStockLockAdapter(inventoryURL string) *HTTPStockLockAdapter {
	return &HTTPStockLockAdapter{
		inventoryURL: inventoryURL,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *HTTPStockLockAdapter) LockStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) ([]string, error) {
	lockKeys := make([]string, 0)
	lockKeyPrefix := fmt.Sprintf("lock-%s-%d", orderID, time.Now().Unix())

	for skuID, qty := range skuQtys {
		lockKey := fmt.Sprintf("%s-%s", lockKeyPrefix, skuID)
		body := map[string]interface{}{
			"sku_id":       skuID,
			"order_id":     orderID,
			"warehouse_id": warehouseID,
			"quantity":     qty,
			"lock_key":     lockKey,
		}
		if err := a.postJSON(ctx, a.inventoryURL+"/api/v1/inventory/lock", body); err != nil {
			return nil, fmt.Errorf("锁定库存失败 sku=%s: %w", skuID, err)
		}
		lockKeys = append(lockKeys, lockKey)
	}
	return lockKeys, nil
}

func (a *HTTPStockLockAdapter) postJSON(ctx context.Context, url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set(headerContentType, mimeJSON)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

// HTTPOutboundCreatorAdapter 通过 HTTP 调用 Warehouse 服务创建出库单
type HTTPOutboundCreatorAdapter struct {
	warehouseURL string
	client       *http.Client
}

func NewHTTPOutboundCreatorAdapter(warehouseURL string) *HTTPOutboundCreatorAdapter {
	return &HTTPOutboundCreatorAdapter{
		warehouseURL: warehouseURL,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *HTTPOutboundCreatorAdapter) CreateOutbound(ctx context.Context, tenantID, orderID, orderNo, warehouseID string, items []OrderItemData) (string, error) {
	itemPayload := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		itemPayload = append(itemPayload, map[string]interface{}{
			"sku_id":   it.SKUID,
			"sku_code": it.SKUCode,
			"sku_name": it.SKUName,
			"quantity": it.Qty,
		})
	}
	body := map[string]interface{}{
		"order_id":     orderID,
		"order_no":     orderNo,
		"warehouse_id": warehouseID,
		"items":        itemPayload,
	}
	_ = tenantID

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.warehouseURL+"/api/v1/warehouse/outbounds", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set(headerContentType, mimeJSON)
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用仓库服务失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("创建出库单失败: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data.ID, nil
}

// HTTPStockDeductAdapter 通过 HTTP 调用 Inventory 服务按订单扣减库存
type HTTPStockDeductAdapter struct {
	inventoryURL string
	client       *http.Client
}

func NewHTTPStockDeductAdapter(inventoryURL string) *HTTPStockDeductAdapter {
	return &HTTPStockDeductAdapter{
		inventoryURL: inventoryURL,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *HTTPStockDeductAdapter) DeductStock(ctx context.Context, orderID, warehouseID string, skuQtys map[string]int) error {
	_ = warehouseID
	_ = skuQtys
	body := map[string]interface{}{"order_id": orderID}
	return a.postJSON(ctx, a.inventoryURL+"/api/v1/inventory/deduct-by-order", body)
}

func (a *HTTPStockDeductAdapter) postJSON(ctx context.Context, url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set(headerContentType, mimeJSON)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

// HTTPInboundHandlerAdapter 通过 HTTP 调用 Inventory 服务增加库存（采购入库）
type HTTPInboundHandlerAdapter struct {
	inventoryURL string
	client       *http.Client
}

func NewHTTPInboundHandlerAdapter(inventoryURL string) *HTTPInboundHandlerAdapter {
	return &HTTPInboundHandlerAdapter{
		inventoryURL: inventoryURL,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *HTTPInboundHandlerAdapter) ReceiveInbound(ctx context.Context, tenantID, purchaseID, warehouseID, supplierID string, items []OrderItemData) (string, error) {
	inboundID := fmt.Sprintf("IN-%s-%d", purchaseID, time.Now().Unix())
	for _, item := range items {
		body := map[string]interface{}{
			"sku_id":       item.SKUID,
			"warehouse_id": warehouseID,
			"quantity":     item.Qty,
			"inbound_id":   inboundID,
			"idempotency_key": fmt.Sprintf("inbound-%s-%s", inboundID, item.SKUID),
		}
		_ = tenantID
		_ = supplierID
		data, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("序列化入库请求失败 sku=%s: %w", item.SKUID, err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.inventoryURL+"/api/v1/inventory/increase", bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("创建入库请求失败 sku=%s: %w", item.SKUID, err)
		}
		req.Header.Set(headerContentType, mimeJSON)
		resp, err := a.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("调用库存服务增加接口失败 sku=%s: %w", item.SKUID, err)
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			return "", fmt.Errorf("增加库存失败 sku=%s: HTTP %d", item.SKUID, resp.StatusCode)
		}
	}
	return inboundID, nil
}
