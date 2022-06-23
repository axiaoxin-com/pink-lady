// 业务错误码定义

package response

import (
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/chai2010/gettext-go"
)

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
	CodeSuccess       = goutils.NewErrCode(success, gettext.Gettext("成功"))
	CodeFailure       = goutils.NewErrCode(failure, gettext.Gettext("失败"))
	CodeInvalidParam  = goutils.NewErrCode(invalidParam, gettext.Gettext("参数错误"))
	CodeNotFound      = goutils.NewErrCode(notFound, gettext.Gettext("未找到"))
	CodeInternalError = goutils.NewErrCode(unknownError, gettext.Gettext("未知错误"))
)

// IsInvalidParamError 判断错误信息中是否包含:参数错误
func IsInvalidParamError(err error) bool {
	return strings.Contains(err.Error(), "Invalid Param")
}
