package repository

import "time"

type NotificationModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;index"`
	UserID    string    `gorm:"column:user_id;index"`
	Title     string    `gorm:"column:title"`
	Content   string    `gorm:"column:content"`
	Type      string    `gorm:"column:type"`
	Read      bool      `gorm:"column:read"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
func (NotificationModel) TableName() string { return "notifications" }
