package pagination

import (
	"math"

	"gorm.io/gorm"
)

// PageRequest 分页请求参数
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// DefaultPageRequest 默认分页参数
func DefaultPageRequest() PageRequest {
	return PageRequest{
		Page:     1,
		PageSize: 20,
	}
}

// Offset 计算偏移量
func (p PageRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 获取限制数
func (p PageRequest) Limit() int {
	return p.PageSize
}

// Paginate GORM 分页查询封装
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// TotalPages 计算总页数
func TotalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
