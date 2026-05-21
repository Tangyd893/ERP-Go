package domain

import "time"

// Notification 通知实体
type Notification struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // info, warning, success, error
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}
