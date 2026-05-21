package errors

import "fmt"

// ErrorCode 统一错误码
type ErrorCode int

const (
	// 成功
	CodeSuccess ErrorCode = 0

	// 通用错误 10000-19999
	CodeInternalError     ErrorCode = 10000
	CodeInvalidParameter  ErrorCode = 10001
	CodeNotFound          ErrorCode = 10002
	CodeAlreadyExists     ErrorCode = 10003
	CodeServiceUnavailable ErrorCode = 10004
	CodeTimeout           ErrorCode = 10005
	CodeRateLimited       ErrorCode = 10006
	CodeForbidden         ErrorCode = 10007

	// 认证权限 20000-29999
	CodeUnauthorized      ErrorCode = 20000
	CodeTokenExpired      ErrorCode = 20001
	CodeTokenInvalid      ErrorCode = 20002
	CodePermissionDenied  ErrorCode = 20003
	CodeLoginFailed       ErrorCode = 20004
	CodeUserDisabled      ErrorCode = 20005

	// 租户 30000-39999
	CodeTenantNotFound    ErrorCode = 30000
	CodeTenantDisabled    ErrorCode = 30001
	CodeTenantQuotaExceeded ErrorCode = 30002

	// 商品 40000-49999
	CodeSKUNotFound       ErrorCode = 40000
	CodeSKUAlreadyExists  ErrorCode = 40001
	CodeSKUMappingNotFound ErrorCode = 40002

	// 订单 50000-59999
	CodeOrderNotFound     ErrorCode = 50000
	CodeOrderStatusInvalid ErrorCode = 50001
	CodeOrderDuplicate    ErrorCode = 50002
	CodeOrderAuditFailed  ErrorCode = 50003

	// 库存 60000-69999
	CodeInsufficientStock ErrorCode = 60000
	CodeStockLockFailed   ErrorCode = 60001
	CodeStockReleaseFailed ErrorCode = 60002
	CodeStockDeductFailed ErrorCode = 60003
	CodeStockIdempotencyConflict ErrorCode = 60004

	// 仓储 70000-79999
	CodeWarehouseNotFound ErrorCode = 70000
	CodeLocationFull      ErrorCode = 70001
	CodePickTaskNotFound  ErrorCode = 70002
	CodeOutboundFailed    ErrorCode = 70003

	// 物流 80000-89999
	CodeCarrierNotFound   ErrorCode = 80000
	CodeLabelCreationFailed ErrorCode = 80001
	CodeTrackingNotFound  ErrorCode = 80002

	// 财务 90000-99999
	CodeSettlementNotFound ErrorCode = 90000
	CodePaymentFailed     ErrorCode = 90001
	CodeCurrencyNotSupported ErrorCode = 90002
	CodeProfitCalculationError ErrorCode = 90003
)

// 错误码对应的中文消息
var codeMessages = map[ErrorCode]string{
	CodeSuccess:            "操作成功",
	CodeInternalError:      "系统内部错误",
	CodeInvalidParameter:   "参数无效",
	CodeNotFound:           "资源未找到",
	CodeAlreadyExists:      "资源已存在",
	CodeServiceUnavailable: "服务暂不可用",
	CodeTimeout:            "请求超时",
	CodeRateLimited:        "请求过于频繁",
	CodeForbidden:          "禁止访问",

	CodeUnauthorized:      "未授权",
	CodeTokenExpired:      "令牌已过期",
	CodeTokenInvalid:      "令牌无效",
	CodePermissionDenied:  "权限不足",
	CodeLoginFailed:       "登录失败",
	CodeUserDisabled:      "用户已禁用",

	CodeTenantNotFound:     "租户未找到",
	CodeTenantDisabled:     "租户已禁用",
	CodeTenantQuotaExceeded: "租户配额已用尽",

	CodeSKUNotFound:       "SKU未找到",
	CodeSKUAlreadyExists:  "SKU已存在",
	CodeSKUMappingNotFound: "SKU映射未找到",

	CodeOrderNotFound:      "订单未找到",
	CodeOrderStatusInvalid: "订单状态无效",
	CodeOrderDuplicate:     "订单重复",
	CodeOrderAuditFailed:   "订单审核失败",

	CodeInsufficientStock:       "库存不足",
	CodeStockLockFailed:         "库存锁定失败",
	CodeStockReleaseFailed:      "库存释放失败",
	CodeStockDeductFailed:       "库存扣减失败",
	CodeStockIdempotencyConflict: "库存幂等冲突",

	CodeWarehouseNotFound: "仓库未找到",
	CodeLocationFull:      "库位已满",
	CodePickTaskNotFound:  "拣货任务未找到",
	CodeOutboundFailed:    "出库失败",

	CodeCarrierNotFound:    "物流商未找到",
	CodeLabelCreationFailed: "面单创建失败",
	CodeTrackingNotFound:   "物流轨迹未找到",

	CodeSettlementNotFound:    "结算单未找到",
	CodePaymentFailed:         "付款失败",
	CodeCurrencyNotSupported:  "不支持该币种",
	CodeProfitCalculationError: "利润计算错误",
}

// Message 返回错误码对应的中文消息
func (c ErrorCode) Message() string {
	if msg, ok := codeMessages[c]; ok {
		return msg
	}
	return "未知错误"
}

// Int 返回错误码的整数值
func (c ErrorCode) Int() int {
	return int(c)
}

// BusinessError 业务错误
type BusinessError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"-"`
}

func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// NewBusinessError 创建业务错误
func NewBusinessError(code ErrorCode, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// NewBusinessErrorf 创建格式化业务错误
func NewBusinessErrorf(code ErrorCode, format string, args ...interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// WrapError 包装原始错误为业务错误
func WrapError(code ErrorCode, cause error) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: code.Message(),
		Cause:   cause,
	}
}

// WrapErrorWithMessage 包装原始错误，自定义消息
func WrapErrorWithMessage(code ErrorCode, message string, cause error) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
