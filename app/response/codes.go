// 定义返回码

package response

// codes
const (
	Success = iota
	Failure
	InvalidParam
	NotFound
	InternalError
)

// RetCodes
var (
	RCSuccess       = NewRetCode(Success, "成功")
	RCFailure       = NewRetCode(Failure, "失败")
	RCInvalidParam  = NewRetCode(InvalidParam, "参数错误")
	RCNotFound      = NewRetCode(NotFound, "资源不存在")
	RCInternalError = NewRetCode(InternalError, "内部错误")
)
