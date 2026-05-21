package response

import (
	"net/http"

	"github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页响应数据
type PageData struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.CodeSuccess),
		Message: "操作成功",
		Data:    data,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.CodeSuccess),
		Message: "操作成功",
		Data: PageData{
			List:       list,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, code errors.ErrorCode, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    int(code),
		Message: message,
	})
}

// BusinessError 业务错误响应
func BusinessError(c *gin.Context, err *errors.BusinessError) {
	httpStatus := http.StatusInternalServerError
	switch {
	case err.Code == errors.CodeInvalidParameter:
		httpStatus = http.StatusBadRequest
	case err.Code == errors.CodeUnauthorized,
		err.Code == errors.CodeTokenExpired,
		err.Code == errors.CodeTokenInvalid,
		err.Code == errors.CodeLoginFailed:
		httpStatus = http.StatusUnauthorized
	case err.Code == errors.CodeForbidden,
		err.Code == errors.CodePermissionDenied:
		httpStatus = http.StatusForbidden
	case err.Code == errors.CodeNotFound,
		err.Code == errors.CodeSKUNotFound,
		err.Code == errors.CodeOrderNotFound,
		err.Code == errors.CodeWarehouseNotFound:
		httpStatus = http.StatusNotFound
	case err.Code == errors.CodeAlreadyExists,
		err.Code == errors.CodeOrderDuplicate,
		err.Code == errors.CodeSKUAlreadyExists:
		httpStatus = http.StatusConflict
	case err.Code == errors.CodeRateLimited:
		httpStatus = http.StatusTooManyRequests
	}
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    int(err.Code),
		Message: err.Message,
	})
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, code errors.ErrorCode, message string) {
	Error(c, http.StatusInternalServerError, code, message)
}
