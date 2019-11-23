// Package retcode define all api return code at here
package retcode

// RetCode has code and message field
type RetCode struct {
	code    int
	message string
}

// Decode return RetCode private field to protected others modify the code and message global
func (rc *RetCode) Decode() (int, string) {
	return rc.code, rc.message
}

// NewRetCode 创建返回码
func NewRetCode(code int, msg string) *RetCode {
	return &RetCode{code: code, message: msg}
}

// define your return code at here
var (
	Success       = NewRetCode(success, "成功")
	Failure       = NewRetCode(failure, "失败")
	UnknownError  = NewRetCode(unknownError, "未知错误")
	InvalidParams = NewRetCode(invalidParams, "无效参数")
	RouteNotFound = NewRetCode(notFound, "路由不存在")
	InternalError = NewRetCode(internalError, "内部错误")
)
