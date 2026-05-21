package domain

import "time"

// ResourceType 资源类型
type ResourceType string

const (
	ResourceAny      ResourceType = "*"
	ResourceMenu     ResourceType = "menu"
	ResourceButton   ResourceType = "button"
	ResourceAPI      ResourceType = "api"
	ResourceData     ResourceType = "data"
)

// Permission 权限实体
type Permission struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Code         string       `json:"code"`
	Description  string       `json:"description"`
	ResourceType ResourceType `json:"resource_type"`
	Action       string       `json:"action"`
	ParentID     string       `json:"parent_id,omitempty"`
	SortOrder    int          `json:"sort_order"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// IsMenu 是否为菜单权限
func (p *Permission) IsMenu() bool {
	return p.ResourceType == ResourceMenu
}

// IsButton 是否为按钮权限
func (p *Permission) IsButton() bool {
	return p.ResourceType == ResourceButton
}
