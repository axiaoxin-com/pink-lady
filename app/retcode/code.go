// 实际json对象中的code定义

package retcode

// 删除已有的iota常量时记得用_占位
const (
	success = iota
	failure
	unknownError
	invalidParams
	notFound
	internalError
)
