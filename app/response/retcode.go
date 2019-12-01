// RetCode 结构体及其方法定义，实现error接口

package response

import "strings"

// RetCode 返回码结构体
type RetCode struct {
	code interface{}
	msg  string
	errs []error
}

// NewRetCode 创建返回码,code可以是任意类型
func NewRetCode(code interface{}, msg string, errs ...error) *RetCode {
	return &RetCode{
		code: code,
		msg:  msg,
		errs: errs,
	}
}

// Decode 返回结构体里面的字段
func (c *RetCode) Decode() (interface{}, string, []error) {
	return c.code, c.msg, c.errs
}

// Error 拼接所有err信息并返回
func (c *RetCode) Error() string {
	errs := []string{}
	for _, err := range c.errs {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, "; ")
}

// Code 返回code
func (c *RetCode) Code() interface{} {
	return c.code
}

// Msg 返回msg描述
func (c *RetCode) Msg() string {
	return c.msg
}

// Errs 返回error列表
func (c *RetCode) Errs() []error {
	return c.errs
}

// SetMsg 更新msg字段返回新的对象避免覆盖原始对象
func (c *RetCode) SetMsg(msg string) *RetCode {
	nc := *c
	nc.msg = msg
	return &nc
}

// AppendError 添加err到errs中
func (c *RetCode) AppendError(errs ...error) *RetCode {
	nc := *c
	nc.errs = append(nc.errs, errs...)
	return &nc
}
