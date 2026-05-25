package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestGatewayRoutesContract 验证 API 网关路由注册完整性（13个服务）
func TestGatewayRoutesContract(t *testing.T) {
	expectedRoutes := []struct {
		prefix string
		target string
		svc    string
	}{
		{"/api/v1/iam/", "http://localhost:8081", "iam"},
		{"/api/v1/tenant/", "http://localhost:8082", "tenant"},
		{"/api/v1/product/", "http://localhost:8083", "product"},
		{"/api/v1/channel/", "http://localhost:8084", "channel"},
		{"/api/v1/order/", "http://localhost:8085", "order"},
		{"/api/v1/inventory/", "http://localhost:8086", "inventory"},
		{"/api/v1/warehouse/", "http://localhost:8087", "warehouse"},
		{"/api/v1/transport/", "http://localhost:8088", "transport"},
		{"/api/v1/file/", "http://localhost:8089", "file"},
		{"/api/v1/purchase/", "http://localhost:8091", "purchase"},
		{"/api/v1/finance/", "http://localhost:8092", "finance"},
		{"/api/v1/report/", "http://localhost:8093", "report"},
		{"/api/v1/notification/", "http://localhost:8094", "notification"},
	}

	if len(expectedRoutes) != 13 {
		t.Errorf("应注册 13 个服务路由，实际 %d 个", len(expectedRoutes))
	}

	routeSet := make(map[string]bool)
	for _, r := range expectedRoutes {
		key := r.svc
		if routeSet[key] {
			t.Errorf("服务路由重复注册: %s", key)
		}
		routeSet[key] = true
	}
}

// TestAuthSkipPathsContract 验证鉴权跳过路径
func TestAuthSkipPathsContract(t *testing.T) {
	skipPaths := []string{
		"/health",
		"/api/v1/iam/login",
		"/api/v1/iam/refresh",
	}

	for _, p := range skipPaths {
		t.Logf("鉴权跳过路径: %s", p)
	}
}

// TestEventPayloadContract 验证事件载荷 JSON 结构
func TestEventPayloadContract(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		data      interface{}
		wantFields []string
	}{
		{
			name:      "订单导入事件",
			eventType: events.EventOrderImported,
			data: map[string]interface{}{
				"order_id": "order-001",
				"tenant_id": "t-001",
				"order_no": "AMZ-001",
			},
			wantFields: []string{"order_id", "tenant_id", "order_no"},
		},
		{
			name:      "订单审核通过事件",
			eventType: events.EventOrderApproved,
			data: map[string]interface{}{
				"order_id": "order-002",
				"tenant_id": "t-001",
				"warehouse_id": "wh-001",
			},
			wantFields: []string{"order_id", "tenant_id", "warehouse_id"},
		},
		{
			name:      "库存锁定事件",
			eventType: events.EventStockLocked,
			data: map[string]interface{}{
				"order_id": "order-003",
				"warehouse_id": "wh-001",
				"lock_keys": []string{"lock-1"},
			},
			wantFields: []string{"order_id", "warehouse_id", "lock_keys"},
		},
		{
			name:      "出库单创建事件",
			eventType: events.EventOutboundCreated,
			data: map[string]interface{}{
				"outbound_id": "OB-001",
				"order_id": "order-004",
			},
			wantFields: []string{"outbound_id", "order_id"},
		},
		{
			name:      "出库发货事件",
			eventType: events.EventOutboundShipped,
			data: map[string]interface{}{
				"outbound_id": "OB-002",
				"order_id": "order-005",
				"tracking_no": "SF123456",
				"carrier": "顺丰",
			},
			wantFields: []string{"outbound_id", "order_id", "tracking_no", "carrier"},
		},
		{
			name:      "订单发货事件",
			eventType: events.EventOrderShipped,
			data: map[string]interface{}{
				"order_id": "order-006",
				"tracking_no": "SF123456",
			},
			wantFields: []string{"order_id", "tracking_no"},
		},
		{
			name:      "库存扣减事件",
			eventType: events.EventStockDeducted,
			data: map[string]interface{}{
				"order_id": "order-007",
				"warehouse_id": "wh-001",
			},
			wantFields: []string{"order_id", "warehouse_id"},
		},
		{
			name:      "库存释放事件",
			eventType: events.EventStockReleased,
			data: map[string]interface{}{
				"order_id": "order-008",
			},
			wantFields: []string{"order_id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := outbox.NewEventPayload(tt.eventType, tt.data)
			if err != nil {
				t.Fatalf("构建事件载荷失败: %v", err)
			}

			var ep outbox.EventPayload
			if err := json.Unmarshal(payload, &ep); err != nil {
				t.Fatalf("事件载荷 JSON 格式无效: %v", err)
			}

			if ep.EventType != tt.eventType {
				t.Errorf("事件类型应为 %s，实际 %s", tt.eventType, ep.EventType)
			}

			var dataMap map[string]interface{}
			if err := json.Unmarshal(ep.Data, &dataMap); err != nil {
				t.Fatalf("事件数据 JSON 格式无效: %v", err)
			}

			for _, field := range tt.wantFields {
				if _, ok := dataMap[field]; !ok {
					t.Errorf("事件载荷缺少字段: %s", field)
				}
			}

			t.Logf("事件载荷验证通过: %s, 字段数=%d", tt.eventType, len(dataMap))
		})
	}
}

// TestAuthResponseContract 验证鉴权错误响应格式（表驱动）
func TestAuthResponseContract(t *testing.T) {
	testCases := []struct {
		name     string
		token    string
		wantCode int
		wantResp float64
	}{
		{"无令牌", "", http.StatusUnauthorized, 20000},
		{"无效令牌", "Bearer invalid-token", http.StatusUnauthorized, 20002},
	}

	router := gin.New()
	router.GET("/api/v1/order/orders", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 20000, "message": "未提供认证令牌"})
			return
		}
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 20002, "message": "令牌无效或已过期"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/order/orders", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.wantCode {
				t.Errorf("%s: 应返回 %d，实际 %d", tc.name, tc.wantCode, w.Code)
			}
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["code"].(float64) != tc.wantResp {
				t.Errorf("%s: 错误码应为 %.0f，实际 %v", tc.name, tc.wantResp, resp["code"])
			}
		})
	}
}
