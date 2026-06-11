package app

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OrderServiceClient 订单服务 HTTP 客户端
type OrderServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrderServiceClient(baseURL string) *OrderServiceClient {
	return &OrderServiceClient{baseURL: baseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

// ImportedOrder 导入订单请求体
type ImportedOrder struct {
	StoreID         string         `json:"store_id"`
	PlatformOrderNo string         `json:"platform_order_no"`
	OrderType       string         `json:"order_type"`
	BuyerName       string         `json:"buyer_name"`
	Currency        string         `json:"currency"`
	TotalAmount     float64        `json:"total_amount"`
	Items           []ImportedItem `json:"items"`
	Address         ImportedAddress `json:"address"`
}

type ImportedItem struct {
	SKUID     string  `json:"sku_id"`
	SKUCode   string  `json:"sku_code"`
	SKUName   string  `json:"sku_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type ImportedAddress struct {
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Country     string `json:"country"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
}

// CreateOrder 调用 order-service 创建订单
func (c *OrderServiceClient) CreateOrder(ctx context.Context, tenantID string, order ImportedOrder) (string, error) {
	body, _ := json.Marshal(order)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/v1/order/orders", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用订单服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("订单服务返回 HTTP %d", resp.StatusCode)
	}
	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Data.ID, nil
}

// ParseCSVOrders 解析 CSV 文件为订单列表
// CSV 格式: store_id,platform_order_no,buyer_name,currency,total_amount,sku_code,sku_name,quantity,unit_price,contact_name,phone,country,city,street
func ParseCSVOrders(reader io.Reader) ([]ImportedOrder, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV 解析失败: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV 文件为空或缺少数据行")
	}

	// 按第一列（store_id + platform_order_no）分组
	orderMap := make(map[string]*ImportedOrder)
	orderKeys := make([]string, 0)

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 14 {
			continue // skip malformed rows
		}
		storeID := strings.TrimSpace(row[0])
		platformNo := strings.TrimSpace(row[1])
		key := storeID + "::" + platformNo

		if _, exists := orderMap[key]; !exists {
			amount, _ := parseFloat(row[4])
			orderMap[key] = &ImportedOrder{
				StoreID:         storeID,
				PlatformOrderNo: platformNo,
				BuyerName:       strings.TrimSpace(row[2]),
				Currency:        strings.TrimSpace(row[3]),
				TotalAmount:     amount,
				OrderType:       "normal",
				Items:           make([]ImportedItem, 0),
				Address: ImportedAddress{
					ContactName: strings.TrimSpace(row[9]),
					Phone:       strings.TrimSpace(row[10]),
					Country:     strings.TrimSpace(row[11]),
					City:        strings.TrimSpace(row[12]),
					StreetLine1: strings.TrimSpace(row[13]),
				},
			}
			orderKeys = append(orderKeys, key)
		}

		qty, _ := parseInt(row[7])
		price, _ := parseFloat(row[8])
		orderMap[key].Items = append(orderMap[key].Items, ImportedItem{
			SKUID:     strings.TrimSpace(row[5]),
			SKUCode:   strings.TrimSpace(row[5]),
			SKUName:   strings.TrimSpace(row[6]),
			Quantity:  qty,
			UnitPrice: price,
		})
	}

	orders := make([]ImportedOrder, len(orderKeys))
	for i, key := range orderKeys {
		orders[i] = *orderMap[key]
	}
	return orders, nil
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%f", &f)
	return f, err
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

// ── Amazon SP-API 适配器接口 ──────────────────────────

// PlatformAdapter 平台适配器接口（Amazon / Shopify / 等）
type PlatformAdapter interface {
	// FetchOrders 拉取平台订单
	FetchOrders(ctx context.Context, storeID string, since time.Time) ([]ImportedOrder, error)
	// PushTracking 回传物流轨迹
	PushTracking(ctx context.Context, storeID, orderNo, trackingNo, carrierCode string) error
	// PlatformCode 返回平台代码
	PlatformCode() string
}

// MockAmazonAdapter Mock Amazon SP-API 适配器（开发用）
type MockAmazonAdapter struct {
	platformCode string
}

func NewMockAmazonAdapter() *MockAmazonAdapter {
	return &MockAmazonAdapter{platformCode: "amazon"}
}

func (m *MockAmazonAdapter) PlatformCode() string { return m.platformCode }

func (m *MockAmazonAdapter) FetchOrders(ctx context.Context, storeID string, since time.Time) ([]ImportedOrder, error) {
	// 返回 Mock 订单数据
	return []ImportedOrder{
		{
			StoreID: storeID, PlatformOrderNo: fmt.Sprintf("AMZ-%d", time.Now().Unix()),
			OrderType: "normal", BuyerName: "Amazon Buyer", Currency: "USD", TotalAmount: 29.99,
			Items: []ImportedItem{
				{SKUID: "sku-001", SKUCode: "A001", SKUName: "Mock商品", Quantity: 1, UnitPrice: 29.99},
			},
			Address: ImportedAddress{
				ContactName: "John Doe", Phone: "1234567890",
				Country: "US", City: "Seattle", StreetLine1: "123 Amazon Way",
			},
		},
	}, nil
}

func (m *MockAmazonAdapter) PushTracking(ctx context.Context, storeID, orderNo, trackingNo, carrierCode string) error {
	// Mock 回传：记录日志
	return nil
}
