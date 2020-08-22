// 定义返回码

package response

import "github.com/axiaoxin-com/goutils"

// 错误码中的 code 定义
const (
	Success = iota
	Failure
	InvalidParam
	NotFound
	InternalError
)

// 错误码对象定义
var (
	CodeSuccess       = goutils.NewErrCode(Success, "成功")
	CodeFailure       = goutils.NewErrCode(Failure, "失败")
	CodeInvalidParam  = goutils.NewErrCode(InvalidParam, "参数错误")
	CodeNotFound      = goutils.NewErrCode(NotFound, "资源不存在")
	CodeInternalError = goutils.NewErrCode(InternalError, "内部错误")
)
