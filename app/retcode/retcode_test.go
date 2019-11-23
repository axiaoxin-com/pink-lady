package retcode

import "testing"

func TestRetCode(t *testing.T) {
	NewCode := &RetCode{code: 999999, message: "xxx"}
	code, msg := NewCode.Decode()
	if code != 999999 && msg != "xxx" {
		t.Error("Retcode decode error")
	}
}
func TestCodeValue(t *testing.T) {
	if code, _ := Success.Decode(); code != success {
		t.Error("success code not equal")
	}
	if code, _ := Failure.Decode(); code != failure {
		t.Error("failure code not equal")
	}
	if code, _ := UnknownError.Decode(); code != unknownError {
		t.Error("unknownError code not equal")
	}
	if code, _ := InvalidParams.Decode(); code != invalidParams {
		t.Error("invalidParams code not equal")
	}
	if code, _ := RouteNotFound.Decode(); code != notFound {
		t.Error("notFound code not equal")
	}
	if code, _ := InternalError.Decode(); code != internalError {
		t.Error("internalError code not equal")
	}
}
