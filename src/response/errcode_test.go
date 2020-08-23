package response

import (
	"testing"
)

func TestErrCode(t *testing.T) {
	if CodeSuccess.Code() != success {
		t.Error("wrong success code")
	}
}
