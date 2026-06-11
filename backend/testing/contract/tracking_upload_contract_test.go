package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestTrackingUploadContract 验证 transport → channel 发货回传契约
func TestTrackingUploadContract(t *testing.T) {
	var capturedStoreID, capturedTrackingNo string

	channelSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/v1/channel/tracking-upload" {
			var req struct {
				StoreID     string `json:"store_id"`
				TrackingNo  string `json:"tracking_no"`
				CarrierCode string `json:"carrier_code"`
				OrderNo     string `json:"order_no"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			if req.StoreID == "" || req.TrackingNo == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			capturedStoreID = req.StoreID
			capturedTrackingNo = req.TrackingNo
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 0,
				"data": map[string]interface{}{
					"id":        "TU-001",
					"task_type": "tracking_upload",
					"status":    "pending",
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer channelSrv.Close()

	// 正确请求
	reqBody := map[string]string{
		"store_id":     "store-001",
		"tracking_no":  "TN-001",
		"carrier_code": "YTO",
		"order_no":     "SO-001",
	}
	body, _ := json.Marshal(reqBody)
	resp, err := http.Post(channelSrv.URL+"/api/v1/channel/tracking-upload", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("契约请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("契约失败: HTTP %d", resp.StatusCode)
	}
	if capturedStoreID != "store-001" {
		t.Errorf("store_id 应为 store-001，实际: %s", capturedStoreID)
	}
	if capturedTrackingNo != "TN-001" {
		t.Errorf("tracking_no 应为 TN-001，实际: %s", capturedTrackingNo)
	}
	t.Log("✓ 发货回传契约通过")

	// 缺少必填字段应拒绝
	badBody, _ := json.Marshal(map[string]string{"store_id": "s1"})
	resp2, _ := http.Post(channelSrv.URL+"/api/v1/channel/tracking-upload", "application/json", bytes.NewReader(badBody))
	if resp2 != nil {
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusBadRequest {
			t.Error("缺少 tracking_no 应返回 400")
		}
	}
	t.Log("✓ 缺少必填字段正确拒绝")
}
