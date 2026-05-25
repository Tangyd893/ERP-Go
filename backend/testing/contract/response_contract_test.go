package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tangyd893/ERP-Go/backend/shared/events"
	"github.com/Tangyd893/ERP-Go/backend/shared/outbox"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponseFormat(t *testing.T) {
	engine := gin.New()
	engine.GET("/test-success", func(c *gin.Context) {
		response.Success(c, gin.H{"id": "1"})
	})
	engine.GET("/test-error", func(c *gin.Context) {
		response.Error(c, http.StatusBadRequest, 10001, "参数错误")
	})
	engine.GET("/test-page", func(c *gin.Context) {
		response.PageSuccess(c, []string{"a", "b"}, 2, 1, 20)
	})

	t.Run("成功响应格式", func(t *testing.T) {
		w := request(t, engine, "GET", "/test-success")
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		if body["code"].(float64) != 0 {
			t.Errorf("成功 code 应为 0，实际 %v", body["code"])
		}
		if _, ok := body["data"]; !ok {
			t.Error("成功响应应包含 data 字段")
		}
	})

	t.Run("错误响应格式", func(t *testing.T) {
		w := request(t, engine, "GET", "/test-error")
		if w.Code != http.StatusBadRequest {
			t.Errorf("应返回 400，实际 %d", w.Code)
		}
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		if body["code"].(float64) != 10001 {
			t.Errorf("错误码不匹配")
		}
	})

	t.Run("分页响应格式", func(t *testing.T) {
		w := request(t, engine, "GET", "/test-page")
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		data := body["data"].(map[string]interface{})
		required := []string{"list", "total", "page", "page_size"}
		for _, f := range required {
			if _, ok := data[f]; !ok {
				t.Errorf("分页响应缺少字段: %s", f)
			}
		}
	})
}

func TestEventPayloadSnapshot(t *testing.T) {
	events_ := map[string]interface{}{
		events.EventOrderApproved: map[string]interface{}{
			"order_id": "order-1", "tenant_id": "default", "store_id": "st-1",
			"order_no": "SO-001", "warehouse_id": "wh-001",
			"items": []map[string]interface{}{{"sku_id": "sku-1", "qty": 2}},
		},
		events.EventOutboundShipped: map[string]interface{}{
			"outbound_id": "OB-1", "order_id": "order-1", "tenant_id": "default",
			"warehouse_id": "wh-001", "tracking_no": "TN001", "carrier": "SF",
			"items": []map[string]interface{}{{"sku_id": "sku-1", "qty": 2}},
		},
		events.EventStockLocked: map[string]interface{}{
			"order_id": "order-1", "warehouse_id": "wh-001",
			"lock_keys": []string{"lock-order-1-sku-1"},
		},
		events.EventOutboundCreated: map[string]interface{}{
			"outbound_id": "OB-1", "order_id": "order-1", "order_no": "SO-001",
			"items": []map[string]interface{}{{"sku_id": "sku-1", "qty": 2}},
		},
		events.EventSettlementImported: map[string]interface{}{
			"purchase_id": "PO-1", "inbound_id": "IN-1", "tenant_id": "default",
		},
	}

	for eventType, data := range events_ {
		t.Run(eventType, func(t *testing.T) {
			payload, err := outbox.NewEventPayload(eventType, data)
			if err != nil {
				t.Fatalf("构建载荷失败: %v", err)
			}
			var ep outbox.EventPayload
			if err := json.Unmarshal(payload, &ep); err != nil {
				t.Fatalf("载荷 JSON 无效: %v", err)
			}
			if ep.EventType != eventType {
				t.Errorf("事件类型应为 %s，实际 %s", eventType, ep.EventType)
			}
		})
	}
}

func request(t *testing.T, handler http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}
