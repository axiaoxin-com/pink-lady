package retcode

const (
	success = iota
	failure
	paramsError
)

type RetCode struct {
	Code    int
	Message string
}

var (
	SUCCESS      = &RetCode{Code: success, Message: "success"}
	FAILURE      = &RetCode{Code: failure, Message: "failure"}
	PARAMS_ERROR = &RetCode{Code: paramsError, Message: "params error"}
)
