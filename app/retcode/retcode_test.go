package retcode

import "testing"

func TestRetCode(t *testing.T) {
	NewCode := &RetCode{code: 999999, message: "xxx"}
	code, msg := NewCode.Decode()
	if code != 999999 && msg != "xxx" {
		t.Error("Retcode decode error")
	}
}
