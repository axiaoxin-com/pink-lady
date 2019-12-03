// Package response 提供统一的JSON返回结构，可以通过配置设置具体返回的code字段为int或者string
package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一的返回结构定义
type Response struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// JSON 返回HTTP状态码为200的统一成功结构
func JSON(c *gin.Context, data interface{}) {
	Respond(c, http.StatusOK, data, RCSuccess)
}

// ErrJSON 返回HTTP状态码为200的统一失败结构
func ErrJSON(c *gin.Context, err error, extraMsgs ...interface{}) {
	Respond(c, http.StatusOK, nil, err, extraMsgs...)
}

// ErrJSON400 respond unified JSON structure with 400 http status code
func ErrJSON400(c *gin.Context, extraMsgs ...interface{}) {
	Respond(c, http.StatusBadRequest, nil, RCInvalidParam, extraMsgs...)
}

// ErrJSON404 respond unified JSON structure with 404 http status code
func ErrJSON404(c *gin.Context, extraMsgs ...interface{}) {
	Respond(c, http.StatusNotFound, nil, RCNotFound, extraMsgs...)
}

// ErrJSON500 respond unified JSON structure with 500 http status code
func ErrJSON500(c *gin.Context, extraMsgs ...interface{}) {
	Respond(c, http.StatusInternalServerError, nil, RCInternalError, extraMsgs...)
}

// Respond encapsulates c.JSON
// debug mode respond indented json
func Respond(c *gin.Context, status int, data interface{}, err error, extraMsgs ...interface{}) {
	// 初始化code、msg为失败
	code, msg, _ := RCFailure.Decode()

	if rc, ok := err.(*RetCode); ok {
		// 如果是返回码，正常处理
		code, msg, _ = rc.Decode()
		// 存在errs则将errs信息添加的msg
		if len(rc.Errs()) > 0 {
			msg = fmt.Sprint(msg, " ", rc.Error())
		}
	} else {
		// 支持rc参数直接传error，如果是error，则将error信息添加到msg
		msg = fmt.Sprint(msg, " ", err.Error())
	}

	// 将extraMsgs添加到msg
	if len(extraMsgs) > 0 {
		msg = fmt.Sprint(msg, "; ", extraMsgs)
	}

	resp := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	if gin.Mode() == gin.ReleaseMode {
		c.JSON(status, resp)
	} else {
		c.IndentedJSON(status, resp)
	}
}
