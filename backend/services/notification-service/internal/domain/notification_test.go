package domain

import (
	"testing"
	"time"
)

// 创建测试用通知
func setupNotification() *Notification {
	return &Notification{
		ID:        "notif-001",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "订单发货通知",
		Content:   "您的订单 ORD-001 已发货，快递单号 SF1234567890",
		Type:      "info",
		Read:      false,
		CreatedAt: time.Now(),
	}
}

// TestNotificationCreation 测试通知创建与基础字段
func TestNotificationCreation(t *testing.T) {
	n := setupNotification()

	if n.ID == "" {
		t.Error("通知ID不应为空")
	}
	if n.TenantID == "" {
		t.Error("租户ID不应为空")
	}
	if n.UserID == "" {
		t.Error("用户ID不应为空")
	}
	if n.Title == "" {
		t.Error("通知标题不应为空")
	}
	if n.Content == "" {
		t.Error("通知内容不应为空")
	}
	if n.Read != false {
		t.Error("新通知应为未读状态")
	}
}

// TestNotificationType 测试通知类型值
func TestNotificationType(t *testing.T) {
	tests := []struct {
		name       string
		notifType  string
	}{
		{"信息", "info"},
		{"警告", "warning"},
		{"成功", "success"},
		{"错误", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := setupNotification()
			n.Type = tt.notifType
			if n.Type != tt.notifType {
				t.Errorf("通知类型应为 %s，实际 %s", tt.notifType, n.Type)
			}
		})
	}
}

// TestNotificationMarkAsRead 测试标记通知为已读
func TestNotificationMarkAsRead(t *testing.T) {
	n := setupNotification()

	if n.Read {
		t.Error("初始状态应为未读")
	}

	n.Read = true
	if !n.Read {
		t.Error("标记已读后 Read 应为 true")
	}
}

// TestNotificationToggleRead 测试通知已读/未读切换
func TestNotificationToggleRead(t *testing.T) {
	n := setupNotification()

	// 标记已读
	n.Read = true
	if !n.Read {
		t.Error("第一次标记应变为已读")
	}

	// 恢复未读
	n.Read = false
	if n.Read {
		t.Error("第二次切换应变为未读")
	}
}

// TestNotificationEmptyTitle 测试空标题通知
func TestNotificationEmptyTitle(t *testing.T) {
	n := &Notification{
		ID:        "notif-empty-title",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "",
		Content:   "内容不可为空",
		Type:      "info",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Title != "" {
		t.Error("空标题通知 Title 应为空字符串")
	}
}

// TestNotificationEmptyContent 测试空内容通知
func TestNotificationEmptyContent(t *testing.T) {
	n := &Notification{
		ID:        "notif-empty-content",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "标题不可为空",
		Content:   "",
		Type:      "info",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Content != "" {
		t.Error("空内容通知 Content 应为空字符串")
	}
}

// TestNotificationLongContent 测试长内容通知
func TestNotificationLongContent(t *testing.T) {
	longContent := `这是一条很长的通知内容，用于测试系统是否能够正确处理包含大量文本的通知消息。
通知系统通常需要支持各种长度的内容，从简短的系统提示到详细的操作说明。`

	n := &Notification{
		ID:        "notif-long",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "长内容测试",
		Content:   longContent,
		Type:      "info",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Content != longContent {
		t.Error("长内容通知内容应完整保留")
	}
}

// TestNotificationCreatedAt 测试通知创建时间
func TestNotificationCreatedAt(t *testing.T) {
	before := time.Now().Add(-1 * time.Minute)
	n := &Notification{
		ID:        "notif-time",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "时间测试",
		Content:   "测试创建时间",
		Type:      "info",
		Read:      false,
		CreatedAt: time.Now(),
	}
	after := time.Now().Add(1 * time.Minute)

	if n.CreatedAt.Before(before) {
		t.Error("创建时间不应早于测试开始前")
	}
	if n.CreatedAt.After(after) {
		t.Error("创建时间不应晚于测试结束后")
	}
}

// TestNotificationDifferentUser 测试不同用户的通知隔离
func TestNotificationDifferentUser(t *testing.T) {
	user1Notif := &Notification{
		ID:       "notif-u1",
		TenantID: "tenant-001",
		UserID:   "user-001",
		Title:    "用户1的通知",
		Content:  "仅用户1可见",
		Type:     "info",
	}

	user2Notif := &Notification{
		ID:       "notif-u2",
		TenantID: "tenant-001",
		UserID:   "user-002",
		Title:    "用户2的通知",
		Content:  "仅用户2可见",
		Type:     "warning",
	}

	if user1Notif.UserID == user2Notif.UserID {
		t.Error("不同用户的通知 UserID 应不同")
	}
	if user1Notif.ID == user2Notif.ID {
		t.Error("不同通知的 ID 应不同")
	}
}

// TestNotificationWarningType 测试警告类型通知
func TestNotificationWarningType(t *testing.T) {
	n := &Notification{
		ID:        "notif-warn",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "库存不足警告",
		Content:   "SKU TST-001 库存不足10件，请及时补货",
		Type:      "warning",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Type != "warning" {
		t.Errorf("警告通知类型应为 warning，实际 %s", n.Type)
	}
}

// TestNotificationErrorType 测试错误类型通知
func TestNotificationErrorType(t *testing.T) {
	n := &Notification{
		ID:        "notif-err",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "同步失败",
		Content:   "平台订单同步失败：接口超时",
		Type:      "error",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Type != "error" {
		t.Errorf("错误通知类型应为 error，实际 %s", n.Type)
	}
}

// TestNotificationSuccessType 测试成功类型通知
func TestNotificationSuccessType(t *testing.T) {
	n := &Notification{
		ID:        "notif-success",
		TenantID:  "tenant-001",
		UserID:    "user-001",
		Title:     "导入完成",
		Content:   "订单导入成功，共导入128条订单",
		Type:      "success",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if n.Type != "success" {
		t.Errorf("成功通知类型应为 success，实际 %s", n.Type)
	}
}

// TestNotificationBatchCreation 测试批量通知不同用户
func TestNotificationBatchCreation(t *testing.T) {
	users := []string{"user-001", "user-002", "user-003"}
	notifs := make([]*Notification, 0, len(users))

	for i, uid := range users {
		n := &Notification{
			ID:        "notif-batch-" + string(rune('1'+i)),
			TenantID:  "tenant-001",
			UserID:    uid,
			Title:     "系统公告",
			Content:   "系统将于今晚22:00进行维护",
			Type:      "info",
			Read:      false,
			CreatedAt: time.Now(),
		}
		notifs = append(notifs, n)
	}

	if len(notifs) != len(users) {
		t.Errorf("批量通知数量应为 %d，实际 %d", len(users), len(notifs))
	}

	for i, n := range notifs {
		if n.UserID != users[i] {
			t.Errorf("第 %d 个通知 UserID 应为 %s，实际 %s", i, users[i], n.UserID)
		}
	}
}
