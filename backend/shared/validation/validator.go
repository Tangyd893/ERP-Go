package validation

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate 校验结构体
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// ShouldBindAndValidate 绑定 JSON 并校验
func ShouldBindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return validate.Struct(obj)
}

// ShouldBindQueryAndValidate 绑定查询参数并校验
func ShouldBindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}
	return validate.Struct(obj)
}
