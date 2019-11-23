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

// define your return code at here
var (
	Success       = &RetCode{code: success, message: "成功"}
	Failure       = &RetCode{code: failure, message: "失败"}
	UnknownError  = &RetCode{code: unknownError, message: "未知错误"}
	InvalidParams = &RetCode{code: invalidParams, message: "无效参数"}
	RouteNotFound = &RetCode{code: notFound, message: "路由不存在"}
	InternalError = &RetCode{code: internalError, message: "内部错误"}
)
