package response

import "testing"

func TestRetCode(t *testing.T) {
	NewCode := NewRetCode(999999, "xxx")
	code, msg, _ := NewCode.Decode()
	if code != 999999 || msg != "xxx" {
		t.Error("Retcode invalid")
	}
}
