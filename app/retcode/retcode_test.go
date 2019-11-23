package retcode

import "testing"

func TestRetCode(t *testing.T) {
	NewCode := NewRetCode(999999, "xxx")
	code, msg := NewCode.Decode()
	if code != 999999 || msg != "xxx" {
		t.Error("Retcode invalid")
	}
}
