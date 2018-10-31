// retcode package define all api return code at here
package retcode

// 删除已有的iota常量时记得用_占位
const (
	success = iota
	unknownError
	invalidParams
	notFound
	internalError
)

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
	Success       = &RetCode{code: success, message: "success"}
	UnknownError  = &RetCode{code: unknownError, message: "unknown error"}
	InvalidParams = &RetCode{code: invalidParams, message: "invalid params"}
	RouteNotFound = &RetCode{code: notFound, message: "route not found"}
	InternalError = &RetCode{code: internalError, message: "internal error"}
)
