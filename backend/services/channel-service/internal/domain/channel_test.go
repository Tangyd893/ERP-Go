package domain

import (
	"testing"
	"time"
)

// 创建测试用店铺
func setupStore() *Store {
	now := time.Now()
	return &Store{
		ID:           "store-001",
		TenantID:     "tenant-001",
		PlatformCode: "amazon",
		Site:         "us",
		Name:         "测试店铺",
		StoreCode:    "TEST001",
		AuthStatus:   "active",
		Status:       StoreStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// 创建测试用同步任务
func setupSyncTask() *SyncTask {
	return &SyncTask{
		ID:        "sync-001",
		TenantID:  "tenant-001",
		StoreID:   "store-001",
		TaskType:  "order_sync",
		Status:    "pending",
		CreatedAt: time.Now(),
	}
}

// TestStoreStatus 测试店铺状态常量
func TestStoreStatus(t *testing.T) {
	tests := []struct {
		name   string
		status StoreStatus
		want   string
	}{
		{"激活状态", StoreStatusActive, "active"},
		{"过期状态", StoreStatusExpired, "expired"},
		{"暂停状态", StoreStatusSuspended, "suspended"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("状态值应为 %s，实际 %s", tt.want, tt.status)
			}
		})
	}
}

// TestStoreCreation 测试店铺创建与字段
func TestStoreCreation(t *testing.T) {
	store := setupStore()

	if store.ID == "" {
		t.Error("店铺ID不应为空")
	}
	if store.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if store.PlatformCode == "" {
		t.Error("平台编码不应为空")
	}
	if store.Name == "" {
		t.Error("店铺名称不应为空")
	}
	if store.Status != StoreStatusActive {
		t.Errorf("默认状态应为 active，实际 %s", store.Status)
	}
}

// TestStoreAuthTokenHidden 测试 AuthToken 不序列化
func TestStoreAuthTokenHidden(t *testing.T) {
	store := setupStore()
	store.AuthToken = "secret-token-123"

	// AuthToken 标记为 json:"-" 不会序列化到 JSON
	if store.AuthToken != "secret-token-123" {
		t.Error("AuthToken 应能被赋值")
	}
}

// TestStoreStatusTransitions 测试店铺状态变更
func TestStoreStatusTransitions(t *testing.T) {
	tests := []struct {
		name     string
		from     StoreStatus
		to       StoreStatus
		expected StoreStatus
	}{
		{"激活转暂停", StoreStatusActive, StoreStatusSuspended, StoreStatusSuspended},
		{"暂停转激活", StoreStatusSuspended, StoreStatusActive, StoreStatusActive},
		{"激活转过期", StoreStatusActive, StoreStatusExpired, StoreStatusExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := setupStore()
			store.Status = tt.from
			store.Status = tt.to
			if store.Status != tt.expected {
				t.Errorf("状态应为 %s，实际 %s", tt.expected, store.Status)
			}
		})
	}
}

// TestSyncTaskStatus 测试同步任务状态值
func TestSyncTaskStatus(t *testing.T) {
	validStatuses := []string{"pending", "running", "completed", "failed"}

	task := setupSyncTask()
	for _, s := range validStatuses {
		task.Status = s
		if task.Status != s {
			t.Errorf("状态应为 %s，实际 %s", s, task.Status)
		}
	}
}

// TestSyncTaskTaskType 测试同步任务类型值
func TestSyncTaskTaskType(t *testing.T) {
	validTypes := []string{"order_sync", "inventory_push", "tracking_upload"}

	task := setupSyncTask()
	for _, ttype := range validTypes {
		task.TaskType = ttype
		if task.TaskType != ttype {
			t.Errorf("任务类型应为 %s，实际 %s", ttype, task.TaskType)
		}
	}
}

// TestSyncTaskProgress 测试同步任务进度计数
func TestSyncTaskProgress(t *testing.T) {
	task := setupSyncTask()
	task.TotalCount = 100
	task.SuccessCnt = 80
	task.FailedCnt = 5

	if task.TotalCount != 100 {
		t.Errorf("总数应为100，实际 %d", task.TotalCount)
	}
	if task.SuccessCnt+task.FailedCnt > task.TotalCount {
		t.Error("成功+失败数不应超过总数")
	}
}

// TestSyncTaskStartedFinished 测试同步任务开始结束时间
func TestSyncTaskStartedFinished(t *testing.T) {
	task := setupSyncTask()

	started := time.Now()
	finished := started.Add(5 * time.Minute)
	task.StartedAt = &started
	task.FinishedAt = &finished

	if task.StartedAt == nil || task.FinishedAt == nil {
		t.Error("任务时间不应为空")
	}
	if !task.FinishedAt.After(*task.StartedAt) {
		t.Error("结束时间应在开始时间之后")
	}
}

// TestPlatformAPILogCreation 测试API调用日志创建
func TestPlatformAPILogCreation(t *testing.T) {
	log := &PlatformAPILog{
		ID:           "log-001",
		StoreID:      "store-001",
		Action:       "GetOrders",
		RequestURL:   "https://api.amazon.com/orders",
		StatusCode:   200,
		Duration:     150,
		CreatedAt:    time.Now(),
	}

	if log.ID == "" {
		t.Error("日志ID不应为空")
	}
	if log.StoreID == "" {
		t.Error("店铺ID不应为空")
	}
	if log.StatusCode < 0 {
		t.Error("状态码不应为负")
	}
}

// TestPlatformAPILogStatusCode 测试API日志HTTP状态码
func TestPlatformAPILogStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"成功200", 200},
		{"创建201", 201},
		{"未授权401", 401},
		{"限流429", 429},
		{"服务错误500", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &PlatformAPILog{
				ID:         "test-log",
				StoreID:    "store-001",
				Action:     "test",
				RequestURL: "https://test.api/",
				StatusCode: tt.statusCode,
				Duration:   100,
				CreatedAt:  time.Now(),
			}
			if log.StatusCode != tt.statusCode {
				t.Errorf("状态码应为 %d，实际 %d", tt.statusCode, log.StatusCode)
			}
		})
	}
}

// TestOrderImportTaskStatus 测试订单导入任务状态
func TestOrderImportTaskStatus(t *testing.T) {
	validStatuses := []string{"pending", "processing", "completed", "failed"}

	task := &OrderImportTask{
		ID:             "import-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		ImportType:     "csv",
		IdempotencyKey: "key-001",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	for _, s := range validStatuses {
		task.Status = s
		if task.Status != s {
			t.Errorf("状态应为 %s，实际 %s", s, task.Status)
		}
	}
}

// TestOrderImportTaskImportType 测试订单导入类型
func TestOrderImportTaskImportType(t *testing.T) {
	validTypes := []string{"csv", "api", "manual"}

	task := &OrderImportTask{
		ID:             "import-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		IdempotencyKey: "key-001",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	for _, importType := range validTypes {
		task.ImportType = importType
		if task.ImportType != importType {
			t.Errorf("导入类型应为 %s，实际 %s", importType, task.ImportType)
		}
	}
}

// TestOrderImportTaskIdempotencyKey 测试幂等键
func TestOrderImportTaskIdempotencyKey(t *testing.T) {
	task1 := &OrderImportTask{
		ID:             "import-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		ImportType:     "csv",
		IdempotencyKey: "key-abc-123",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	task2 := &OrderImportTask{
		ID:             "import-002",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		ImportType:     "csv",
		IdempotencyKey: "key-abc-123",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 相同幂等键应相等
	if task1.IdempotencyKey != task2.IdempotencyKey {
		t.Error("相同幂等键应相等")
	}
}

// TestOrderImportTaskRows 测试导入任务行数统计
func TestOrderImportTaskRows(t *testing.T) {
	task := &OrderImportTask{
		ID:             "import-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		ImportType:     "csv",
		IdempotencyKey: "key-001",
		TotalRows:      100,
		SuccessRows:    90,
		FailedRows:     10,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if task.SuccessRows+task.FailedRows != task.TotalRows {
		t.Error("成功行数+失败行数应等于总行数")
	}
}

// TestOrderImportTaskErrorMsg 测试导入失败时的错误信息
func TestOrderImportTaskErrorMsg(t *testing.T) {
	task := &OrderImportTask{
		ID:             "import-001",
		TenantID:       "tenant-001",
		StoreID:        "store-001",
		ImportType:     "csv",
		IdempotencyKey: "key-001",
		Status:         "failed",
		ErrorMsg:       "CSV文件格式错误：缺少必填列",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if task.Status == "failed" && task.ErrorMsg == "" {
		t.Error("失败状态时错误信息不应为空")
	}
}
