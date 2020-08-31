// 业务错误码定义

package response

import "github.com/axiaoxin-com/goutils"

// 错误码中的 code 定义
const (
	failure = iota - 1
	success
	invalidParam
	notFound
	unknownError
)

// 错误码对象定义
var (
	CodeSuccess       = goutils.NewErrCode(success, "成功")
	CodeFailure       = goutils.NewErrCode(failure, "失败")
	CodeInvalidParam  = goutils.NewErrCode(invalidParam, "参数错误")
	CodeNotFound      = goutils.NewErrCode(notFound, "资源不存在")
	CodeInternalError = goutils.NewErrCode(unknownError, "未知错误")
)
